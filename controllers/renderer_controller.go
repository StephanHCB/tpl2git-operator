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
	logger := r.Log.WithValues("renderer", req.NamespacedName)

	// TODO some pointers on how to do proper error handling
	// https://github.com/improbable-eng/etcd-cluster-operator/blob/f84abc6561735814debd67d45bb62d2d2ed8cf4a/controllers/etcdcluster_controller.go#L546

	// obtain the Renderer instance for this reconcile request
	renderer := &tpl2gitv1alpha1.Renderer{}
	if err := r.Get(ctx, req.NamespacedName, renderer); err != nil {
		logger.Error(err, "error during resource get for renderer")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// if you need it:
	// original := renderer.DeepCopy()

	// your logic here

	// if we have status fields, patch the renderer in the cluster to update the status fields

	logger.Info("success updating resource with target_repo_url = " + renderer.Spec.TargetRepoUrl)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RendererReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tpl2gitv1alpha1.Renderer{}).
		Complete(r)
}
