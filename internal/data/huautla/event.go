package huautla

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jsmit257/huautla/types"
)

func (ha *HuautlaAdaptor) PostLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostEvent")
	defer r.Body.Close()

	var e types.Event

	if id, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: event id", err), http.StatusBadRequest, err)
	} else if err = bodyHelper(r, e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(id), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.AddLifecycleEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.created(w, l)
	}
}

func (ha *HuautlaAdaptor) PatchEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchEvent")
	defer r.Body.Close()

	ms.error(w, fmt.Errorf("not implemented"), http.StatusNotImplemented, "events really need to be simpler")

	// var e types.Event

	// if _, err := getUUIDByName("ev_id", w, r, ms); err != nil {
	// 	ms.error(w, fmt.Errorf("%w: event id", err), http.StatusBadRequest, err)
	// } else if body, err := io.ReadAll(r.Body); err != nil {
	// 	ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	// } else if err := json.Unmarshal(body, &e); err != nil {
	// 	ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	// } else if _, err := ha.db.UpdateEvent(r.Context(), e, ms.cid); err != nil {
	// 	ms.error(w, err, http.StatusInternalServerError, "failed to change event")
	// } else {
	// 	ms.empty(w)
	// }
}

func (ha *HuautlaAdaptor) PatchLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchLifecycleEvent")
	defer r.Body.Close()

	var e types.Event

	if lcID, err := getUUIDByName("lc_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: lifecycle id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if _, err := ha.db.ChangeLifecycleEvent(r.Context(), &l, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change event")
	} else {
		ms.send(w, http.StatusOK, l)
	}
}

func (ha *HuautlaAdaptor) DeleteLifecycleEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteLifecycleEvent")

	if lcID, err := getUUIDByName("lc_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: lifecycle id", err), http.StatusBadRequest, err)
	} else if evID, err := getUUIDByName("ev_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: event id", err), http.StatusBadRequest, err)
	} else if l, err := ha.db.SelectLifecycle(r.Context(), types.UUID(lcID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.RemoveLifecycleEvent(r.Context(), &l, types.UUID(evID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove event")
	} else {
		ms.send(w, http.StatusOK, l)
	}
}

func (ha *HuautlaAdaptor) PostGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PostGenerationEvent")
	defer r.Body.Close()

	var e types.Event

	if genID, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(genID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.AddGenerationEvent(r.Context(), &g, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to add event")
	} else {
		ms.send(w, http.StatusCreated, g)
	}
}

func (ha *HuautlaAdaptor) PatchGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "PatchGenerationEvent")
	defer r.Body.Close()

	var e types.Event

	if genID, err := getUUIDByName("id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if body, err := io.ReadAll(r.Body); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't read request body")
	} else if err := json.Unmarshal(body, &e); err != nil {
		ms.error(w, err, http.StatusBadRequest, "couldn't unmarshal request body")
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(genID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if _, err := ha.db.ChangeGenerationEvent(r.Context(), &g, e, ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to change event")
	} else {
		ms.send(w, http.StatusOK, g)
	}
}

func (ha *HuautlaAdaptor) DeleteGenerationEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ms := ha.start(ctx, "DeleteGenerationEvent")

	if gID, err := getUUIDByName("g_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if evID, err := getUUIDByName("ev_id", w, r, ms); err != nil {
		ms.error(w, fmt.Errorf("%w: generation id", err), http.StatusBadRequest, err)
	} else if g, err := ha.db.SelectGeneration(r.Context(), types.UUID(gID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to fetch lifecycle")
	} else if err := ha.db.RemoveGenerationEvent(r.Context(), &g, types.UUID(evID), ms.cid); err != nil {
		ms.error(w, err, http.StatusInternalServerError, "failed to remove lifecycle")
	} else {
		ms.send(w, http.StatusOK, g)
	}
}
