package huautla

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) GetAllVendors(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetAllVendors")
	defer ms.end()

	if vendors, err := ha.db.SelectAllVendors(r.Context(), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch vendors")
	} else {
		ms.send(w, vendors, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) GetVendor(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetVendor")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if vendor, err := ha.db.SelectVendor(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch vendor")
	} else {
		ms.send(w, vendor, http.StatusOK)
	}
}

func (ha *HuautlaAdaptor) PostVendor(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PostVendor")
	defer ms.end()
	defer r.Body.Close()

	var v types.Vendor

	if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &v); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if v, err = ha.db.InsertVendor(r.Context(), v, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to insert vendor")
	}

	ms.send(w, v, http.StatusCreated)
}

func (ha *HuautlaAdaptor) PatchVendor(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("PatchVendor")
	defer ms.end()
	defer r.Body.Close()

	var v types.Vendor

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body") // XXX: better status code??
	} else if err := json.Unmarshal(body, &v); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body") // XXX: better status code??
	} else if err = ha.db.UpdateVendor(r.Context(), types.UUID(id), v, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to update vendor")
	}

	ms.send(w, nil, http.StatusNoContent)
}

func (ha *HuautlaAdaptor) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("DeleteVendor")
	defer ms.end()

	if id := chi.URLParam(r, "id"); id == "" {
		ms.error(w, fmt.Errorf("missing required id parameter"), http.StatusBadRequest, "missing required id parameter")
	} else if id, err := url.QueryUnescape(id); err != nil {
		ms.error(w, fmt.Errorf("malformed id parameter"), http.StatusBadRequest, "malformed id parameter")
	} else if err := ha.db.DeleteVendor(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to delete vendor")
	} else {
		ms.send(w, nil, http.StatusNoContent)
	}
}

func (ha *HuautlaAdaptor) GetVendorReport(w http.ResponseWriter, r *http.Request) {
	ms := ha.start("GetVendorReport")
	defer ms.end()

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch uuid")
	} else if v, err := ha.db.VendorReport(r.Context(), id, ms.cid); errors.Is(err, sql.ErrNoRows) {
		ms.error(w, err, http.StatusBadRequest, "failed to fetch vendor")
	} else if err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch vendor")
	} else {
		ms.send(w, v, http.StatusOK)
	}
}
