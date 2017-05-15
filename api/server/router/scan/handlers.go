package scan

import (
	"encoding/json"
	"net/http"

	"github.com/maliceio/malice/api/server/httputils"
	"github.com/maliceio/malice/api/types/filters"
	scantypes "github.com/maliceio/malice/api/types/scan"
	"golang.org/x/net/context"
)

func (v *scanRouter) getscansList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	scans, warnings, err := v.backend.scans(r.Form.Get("filters"))
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, &scantypes.scansListOKBody{scans: scans, Warnings: warnings})
}

func (v *scanRouter) postscansCreate(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	if err := httputils.CheckForJSON(r); err != nil {
		return err
	}

	var req scantypes.scansCreateBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	scan, err := v.backend.scanCreate(req.Name, req.Driver, req.DriverOpts, req.Labels)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusCreated, scan)
}

func (v *scanRouter) deletescans(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}
	force := httputils.BoolValue(r, "force")
	if err := v.backend.scanRm(vars["name"], force); err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (v *scanRouter) postscansPrune(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := httputils.ParseForm(r); err != nil {
		return err
	}

	pruneFilters, err := filters.FromParam(r.Form.Get("filters"))
	if err != nil {
		return err
	}

	pruneReport, err := v.backend.scansPrune(ctx, pruneFilters)
	if err != nil {
		return err
	}
	return httputils.WriteJSON(w, http.StatusOK, pruneReport)
}
