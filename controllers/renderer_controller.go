/*
Copyright 2021 StephanHCB.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tpl2gitv1alpha1 "github.com/StephanHCB/tpl2git-operator/api/v1alpha1"
)

// RendererReconciler reconciles a Renderer object
type RendererReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tpl2git.stephanhcb.github.io,resources=renderers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tpl2git.stephanhcb.github.io,resources=renderers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tpl2git.stephanhcb.github.io,resources=renderers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Renderer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *RendererReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("renderer", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RendererReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tpl2gitv1alpha1.Renderer{}).
		Complete(r)
}
