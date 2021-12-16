package controllers

import (
	"context"

	appsApi "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	godocApi "github.com/rpflynn22/godoc-operator/internal/api/v1alpha1"
	managed "github.com/rpflynn22/godoc-operator/internal/managed"
)

// RepoReconciler reconciles a Repo object
type RepoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *RepoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("namespace", req.Namespace, "name", req.Name)
	logger.Info("reconciling")

	repo := godocApi.Repo{}
	if err := r.Get(ctx, req.NamespacedName, &repo); err != nil {
		logger.Error(err, "get object")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if isBeingDeleted(&repo) {
		logger.Info("being deleted")
		// delete dependent resources, okay if missing
		// remove finalizer & update
		return ctrl.Result{}, nil
	}

	// add finalizer if not exists

	// create dependent objects in mem
	//   set ownership reference for each dependent object
	deployment := managed.Deployment(&repo)
	if err := controllerutil.SetControllerReference(&repo, deployment, r.Scheme); err != nil {
		logger.Error(err, "deployment set controller reference")
		return ctrl.Result{}, nil
	}

	service := managed.Service(&repo)
	if err := controllerutil.SetControllerReference(&repo, service, r.Scheme); err != nil {
		logger.Error(err, "service set controller reference")
		return ctrl.Result{}, nil
	}

	// look for existing dependent objects
	//   create them if they don't exist
	//   update them if they're misconfigured
	var existingDeployment appsApi.Deployment
	err := r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, &existingDeployment)

	if err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.Create(ctx, deployment); err != nil && !apierrors.IsAlreadyExists(err) {
				logger.Error(err, "create deployment")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "get existing deployment")
			return ctrl.Result{}, err
		}
	}

	var existingService v1.Service
	err = r.Get(ctx, types.NamespacedName{
		Namespace: req.Namespace,
		Name:      req.Name,
	}, &existingService)

	if err != nil {
		if apierrors.IsNotFound(err) {
			if err := r.Create(ctx, service); err != nil && !apierrors.IsAlreadyExists(err) {
				logger.Error(err, "create service")
				return ctrl.Result{}, err
			}
		} else {
			logger.Error(err, "get existing service")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RepoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&godocApi.Repo{}).
		Owns(&appsApi.Deployment{}).
		Owns(&v1.Service{}).
		Complete(r)
}

func isBeingDeleted(repo *godocApi.Repo) bool {
	return !repo.GetDeletionTimestamp().IsZero()
}
