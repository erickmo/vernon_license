package http

import (
	"encoding/json"
	"errors"
	"net/http"

	setupinstall "github.com/flashlab/flasherp-developer-api/internal/command/setup_install"
	getsetupstatus "github.com/flashlab/flasherp-developer-api/internal/query/get_setup_status"
	"github.com/flashlab/flasherp-developer-api/pkg/commandbus"
	"github.com/flashlab/flasherp-developer-api/pkg/querybus"
)

type SetupHandler struct {
	cmdBus        *commandbus.CommandBus
	queryBus      *querybus.QueryBus
	installHandler *setupinstall.Handler
}

func NewSetupHandler(
	cmdBus *commandbus.CommandBus,
	queryBus *querybus.QueryBus,
	installHandler *setupinstall.Handler,
) *SetupHandler {
	return &SetupHandler{
		cmdBus:         cmdBus,
		queryBus:       queryBus,
		installHandler: installHandler,
	}
}

func (h *SetupHandler) Status(w http.ResponseWriter, r *http.Request) {
	result, err := querybus.Dispatch[*getsetupstatus.SetupStatusResult](
		r.Context(), h.queryBus,
		getsetupstatus.GetSetupStatus{},
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "gagal cek status instalasi")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *SetupHandler) Install(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "request tidak valid")
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "name, email, dan password wajib diisi")
		return
	}
	if len(req.Password) < 8 {
		respondError(w, http.StatusBadRequest, "password minimal 8 karakter")
		return
	}

	if err := h.installHandler.Handle(r.Context(), setupinstall.SetupInstall{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}); err != nil {
		if errors.Is(err, setupinstall.ErrAlreadyInstalled) {
			respondError(w, http.StatusConflict, "sistem sudah terinstal")
			return
		}
		respondError(w, http.StatusInternalServerError, "gagal instalasi: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, map[string]string{
		"message": "instalasi berhasil, silakan login",
	})
}
