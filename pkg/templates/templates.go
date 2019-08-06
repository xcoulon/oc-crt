package templates

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	authv1 "github.com/openshift/api/authorization/v1"
	projectv1 "github.com/openshift/api/project/v1"
	octemplateapi "github.com/openshift/api/template"
	templatev1 "github.com/openshift/api/template/v1"
	authv1client "github.com/openshift/client-go/authorization/clientset/versioned/typed/authorization/v1"
	projectv1client "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/openshift/library-go/pkg/template/generator"
	"github.com/openshift/library-go/pkg/template/templateprocessing"
	errs "github.com/pkg/errors"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	errors "k8s.io/apimachinery/pkg/util/errors"
	restclient "k8s.io/client-go/rest"
	kcmdutil "k8s.io/kubernetes/pkg/kubectl/cmd/util"
)

// Apply applies the given template
func Apply(filename string, r io.Reader, projectName, username string, scheme *runtime.Scheme, f kcmdutil.Factory, config *restclient.Config) error {
	infos, err := f.NewBuilder().
		WithScheme(scheme, scheme.PrioritizedVersionsAllGroups()...).
		LocalParam(true).
		Stream(r, filename).
		Do().
		Infos()
	if err != nil {
		return fmt.Errorf("failed to read input object (not a Template?): %v", err)
	}
	if len(infos) > 1 {
		// in order to run validation on the input given to us by a user, we only support the processing
		// of one template in a list. For instance, we want to be able to fail when a user does not give
		// a parameter that the template wants or when they give a parameter the template doesn't need,
		// as this may indicate that they have mis-used `oc process`. This is much less complicated when
		// we process at most one template.
		log.Warnf("%d input templates found, but only the first will be processed\n", len(infos))
	}
	obj, ok := infos[0].Object.(*templatev1.Template)
	if !ok {
		return errs.Errorf("template is not valid")
	}
	if errs := injectUserVars(map[string]string{
		"PROJECT_NAME":    projectName,
		"ADMIN_USER_NAME": username,
	}, obj, false); errs != nil {
		return errors.NewAggregate(errs)
	}
	tmpl, err := processTemplateLocally(obj, scheme)
	if err != nil {
		return err
	}
	// return nil
	return createResources(tmpl.Objects, projectName, config)
}

// injectUserVars injects user specified variables into the Template
func injectUserVars(values map[string]string, t *templatev1.Template, ignoreUnknownParameters bool) []error {
	var errors []error
	for param, val := range values {
		v := templateprocessing.GetParameterByName(t, param)
		if v != nil {
			v.Value = val
			v.Generate = ""
		} else if !ignoreUnknownParameters {
			errors = append(errors, fmt.Errorf("unknown parameter name %q\n", param))
		}
	}
	return errors
}

// processTemplateLocally applies the same logic that a remote call would make but makes no
// connection to the server.
func processTemplateLocally(tpl *templatev1.Template, scheme *runtime.Scheme) (*templatev1.Template, error) {
	processor := templateprocessing.NewProcessor(map[string]generator.Generator{
		"expression": generator.NewExpressionValueGenerator(rand.New(rand.NewSource(time.Now().UnixNano()))),
	})
	if errs := processor.Process(tpl); len(errs) > 0 {
		return nil, kerrors.NewInvalid(octemplateapi.Kind("Template"), tpl.Name, errs)
	}
	var externalResultObj templatev1.Template
	if err := scheme.Convert(tpl, &externalResultObj, nil); err != nil {
		return nil, fmt.Errorf("unable to convert template to external template object: %v", err)
	}
	return &externalResultObj, nil
}

func createResources(objs []runtime.RawExtension, namespace string, config *restclient.Config) error {
	for _, obj := range objs {
		if obj.Object == nil {
			log.Warn("template object is nil")
			continue
		}
		gvk := obj.Object.GetObjectKind().GroupVersionKind()
		log.Infof("processing object of group/kind '%v/%v'", gvk.Version, gvk.Kind)
		switch gvk.Kind {
		case "ProjectRequest":
			projectClient, err := projectv1client.NewForConfig(config)
			if err != nil {
				return errs.Wrapf(err, "failed to init the project client")
			}
			processed := &projectv1.ProjectRequest{}
			err = projectClient.RESTClient().Post().
				Resource("projectrequests").
				Body(obj.Object).
				Do().
				Into(processed)
			if err != nil {
				return errs.Wrapf(err, "failed to create the Project resource")
			}
		case "RoleBinding":
			authClient, err := authv1client.NewForConfig(config)
			if err != nil {
				return errs.Wrapf(err, "failed to init the auth client")
			}
			processed := &authv1.RoleBinding{}
			err = authClient.RESTClient().Post().
				Namespace(namespace).
				Resource("rolebindings").
				Body(obj.Object).
				Do().
				Into(processed)
			if err != nil {
				return errs.Wrapf(err, "failed to create the RoleBinding resource")
			}
		}
	}
	return nil
}
