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

package controllers

import (
	danaiov1 "NamespaceLabel/api/v1"
	"context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/source"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	//Log	logr.Logger
	Scheme    *runtime.Scheme
	NslEvents chan event.GenericEvent
}

//+kubebuilder:rbac:groups=dana.io,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=dana.io,resources=namespaces/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=dana.io,resources=namespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	//l := r.Log.WithValues("namespace", req.NamespacedName)
	l := log.FromContext(ctx)

	type Set map[string]string
	//
	var ns corev1.Namespace
	var labels map[string]string
	var nsll danaiov1.NamespaceLabelList
	if err := r.List(ctx, &nsll, client.InNamespace(req.Name)); err != nil {
		if apierrors.IsNotFound(err) {
			// we'll ignore not-found errors, since we can get them on deleted requests.
			return ctrl.Result{}, nil
		}
		l.Error(err, "unable to get list of nsl")
		return ctrl.Result{}, err
	}

	if err := r.Get(ctx, types.NamespacedName{Name: req.Name}, &ns); err != nil {
		if apierrors.IsNotFound(err) {
			// we'll ignore not-found errors, since we can get them on deleted requests.
			return ctrl.Result{}, nil
		}
		l.Error(err, "unable to fetch namespace")
		return ctrl.Result{}, err
	}

	for _, nsl := range nsll.Items {
		labels = Merge(labels, nsl.Spec.Labels)
	}

	ns.SetLabels(labels)
	if err := r.Update(ctx, &ns); err != nil {
		if apierrors.IsConflict(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		if apierrors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		l.Error(err, "unable to update namespace")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func Merge(labels1, labels2 Set) Set {
	mergedMap := Set{}

	for k, v := range labels1 {
		mergedMap[k] = v
	}
	for k, v := range labels2 {
		mergedMap[k] = v
	}
	return mergedMap
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		// Uncomment the following line adding a pointer to an instance of the controlled resource as an argument
		// For().
		For(&corev1.Namespace{}).
		Watches(&source.Channel{Source: r.NslEvents}, &handler.EnqueueRequestForObject{}).
		Owns(&danaiov1.NamespaceLabel{}).
		Complete(r)
}
