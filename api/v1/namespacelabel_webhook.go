/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

type Set map[string]string

//blacklist for labels
var blacklist = []string{"app", "dana"}

func (ls Set) Has(label string) bool {
	_, exists := ls[label]
	return exists
}

func AreLabelsInBlackList(labels Set, blacklist []string) bool {
	if len(blacklist) == 0 {
		return false
	}

	for _, key := range blacklist {
		if labels.Has(key) {
			return true
		}
	}
	return false
}

// log is for logging in this package.
var namespacelabellog = logf.Log.WithName("namespacelabel-resource")

func (r *NamespaceLabel) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-dana-io-v1-namespacelabel,mutating=false,failurePolicy=fail,sideEffects=None,groups=dana.io,resources=namespacelabels,verbs=create;update,versions=v1,name=vnamespacelabel.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &NamespaceLabel{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *NamespaceLabel) ValidateCreate() error {
	namespacelabellog.Info("validate create", "name", r.Name)

	return r.validateNamespaceLabel()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *NamespaceLabel) ValidateUpdate(old runtime.Object) error {
	namespacelabellog.Info("validate update", "name", r.Name)

	return r.validateNamespaceLabel()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *NamespaceLabel) ValidateDelete() error {
	namespacelabellog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *NamespaceLabel) validateNamespaceLabel() error {
	var allErrs field.ErrorList
	if err := r.validateNamespaceLabelSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "dana.io", Kind: "NamespaceLabel"}, r.Name, allErrs)
}

func (r *NamespaceLabel) validateNamespaceLabelSpec() *field.Error {
	if AreLabelsInBlackList(r.Spec.Labels, blacklist) {
		return field.Invalid(field.NewPath("spec").Child("labels"), r.Spec.Labels, "has forbidden label")
	}
	return nil
}
