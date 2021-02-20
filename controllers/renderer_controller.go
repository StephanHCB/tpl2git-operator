package controllers

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	generatorgit "github.com/StephanHCB/go-generator-git"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tpl2gitv1alpha1 "github.com/StephanHCB/tpl2git-operator/api/v1alpha1"
)

func isUnchanged(spec tpl2gitv1alpha1.RendererSpec, status tpl2gitv1alpha1.RendererStatus) bool {
	if spec.BlueprintRepoUrl != status.CurrentBlueprintRepoUrl {
		return false
	}
	if spec.BlueprintBranch != status.CurrentBlueprintBranch {
		return false
	}
	if spec.BlueprintName != status.CurrentBlueprintName {
		return false
	}
	if spec.TargetRepoUrl != status.CurrentTargetRepoUrl {
		return false
	}
	if spec.TargetBranch != status.CurrentTargetBranch {
		return false
	}
	if spec.TargetBranchForkFrom != status.CurrentTargetBranchForkFrom {
		return false
	}
	if spec.TargetSpecFile != status.CurrentTargetSpecFile {
		return false
	}
	return cmp.Equal(spec.Parameters, status.CurrentParameters, cmpopts.EquateEmpty())
}

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

	// bail out if status and spec are in sync to avoid needless runs
	if isUnchanged(renderer.Spec, renderer.Status) {
		logger.Info("no update needed, success")
		return ctrl.Result{}, nil
	}

	// TODO create Logr integration library for go-autumn-logging
	// for now, just avoid the nil deref
	aulogging.SetupNoLoggerForTesting()

	// business logic

	gen := generatorgit.ThreadsafeInstance()

	if err := gen.CreateTemporaryWorkdir(ctx, "/tmp"); err != nil {
		logger.Error(err, "error during CreateTemporaryWorkdir")
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}

	if err := gen.CloneSourceRepo(ctx, renderer.Spec.BlueprintRepoUrl, renderer.Spec.BlueprintBranch); err != nil {
		logger.Error(err, "error during CloneSourceRepo")
		_ = gen.Cleanup(ctx)
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}

	if err := gen.CloneTargetRepo(ctx, renderer.Spec.TargetRepoUrl, renderer.Spec.TargetBranch, renderer.Spec.TargetBranchForkFrom); err != nil {
		logger.Error(err, "error during CloneTargetRepo")
		_ = gen.Cleanup(ctx)
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}

	if response, err := gen.WriteRenderSpecFile(ctx, renderer.Spec.BlueprintName, "values.txt", renderer.Spec.Parameters); err != nil {
		logger.Error(err, "error(s) during WriteRenderSpecFile")
		for i, e := range response.Errors {
			logger.Error(e, fmt.Sprintf("error %d: %s", i+1, e.Error()))
		}

		_ = gen.Cleanup(ctx)
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}

	if response, err := gen.Generate(ctx); err != nil {
		logger.Error(err, "error(s) during Generate")
		for i, e := range response.Errors {
			logger.Error(e, fmt.Sprintf("error %d: %s", i+1, e.Error()))
		}

		_ = gen.Cleanup(ctx)
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}

	// TODO CommitAndPush
	// - needs name, email, message fields in CRD
	// - needs auth info

	if err := gen.Cleanup(ctx); err != nil {
		logger.Error(err, "error during Cleanup")
	}

	renderer.Status.CurrentBlueprintRepoUrl = renderer.Spec.BlueprintRepoUrl
	renderer.Status.CurrentBlueprintBranch = renderer.Spec.BlueprintBranch
	renderer.Status.CurrentBlueprintName = renderer.Spec.BlueprintName
	renderer.Status.CurrentTargetRepoUrl = renderer.Spec.TargetRepoUrl
	renderer.Status.CurrentTargetBranch = renderer.Spec.TargetBranch
	renderer.Status.CurrentTargetBranchForkFrom = renderer.Spec.TargetBranchForkFrom
	renderer.Status.CurrentTargetSpecFile = renderer.Spec.TargetSpecFile
	renderer.Status.CurrentParameters = renderer.Spec.Parameters

	// update the renderer in the cluster to write back the status
	if err := r.Status().Update(ctx, renderer); err != nil {
		logger.Error(err, "error during status update for renderer")
		// TODO for now, ignore errors and return nil so we do not get continuously rescheduled
		return ctrl.Result{}, nil
	}
	// this updates the base resource, not the status
	//   r.Update(ctx, renderer)
	// thereby creating an event loop

	logger.Info("success updating resource with target_repo_url = " + renderer.Spec.TargetRepoUrl)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RendererReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tpl2gitv1alpha1.Renderer{}).
		Complete(r)
}
