package devices

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type DeviceHandler struct {
	service *Queries
}

func NewDeviceHandler(service *Queries) *DeviceHandler {
	return &DeviceHandler{service: service}
}

func (h *DeviceHandler) CreteDevice(requestWriter http.ResponseWriter, request *http.Request) {
	var device Device
	err := json.NewDecoder(request.Body).Decode(&device)
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.service.CreateDevice(request.Context(), CreateDeviceParams{
		Name:  device.Name,
		Brand: device.Brand,
		Model: device.Model,
		Mac:   device.Mac,
	})
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) GetDeviceByID(requestWriter http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(requestWriter, "id is required", http.StatusBadRequest)
		return
	}
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(requestWriter, "id is invalid", http.StatusBadRequest)
		return
	}
	device, err := h.service.GetDeviceByID(request.Context(), int32(deviceID))
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	deviceDto := NewDeviceDto(device)
	json.NewEncoder(requestWriter).Encode(deviceDto)
}

func (h *DeviceHandler) DeleteDevice(requestWriter http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(requestWriter, "id is required", http.StatusBadRequest)
		return
	}
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(requestWriter, "id is invalid", http.StatusBadRequest)
		return
	}
	_, err = h.service.DeleteDevice(request.Context(), int32(deviceID))
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) UpdateDevice(requestWriter http.ResponseWriter, request *http.Request) {
	var device Device
	err := json.NewDecoder(request.Body).Decode(&device)
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = h.service.UpdateDevice(request.Context(), UpdateDeviceParams{
		ID:    device.ID,
		Name:  device.Name,
		Brand: device.Brand,
		Model: device.Model,
		Mac:   device.Mac,
	})
	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *DeviceHandler) GetDevicePage(requestWriter http.ResponseWriter, request *http.Request) {
	page := request.URL.Query().Get("page")
	if page == "" {
		http.Error(requestWriter, "page is required", http.StatusBadRequest)
		return
	}
	pageNumber, err := strconv.Atoi(page)
	if err != nil {
		http.Error(requestWriter, "page is invalid", http.StatusBadRequest)
		return
	}
	perPage := request.URL.Query().Get("perPage")
	if perPage == "" {
		http.Error(requestWriter, "perPage is required", http.StatusBadRequest)
		return
	}
	perPageNumber, err := strconv.Atoi(perPage)
	if err != nil {
		http.Error(requestWriter, "perPage is invalid", http.StatusBadRequest)
		return
	}

	devices, err := h.service.GetDevicesListPaged(request.Context(), GetDevicesListPagedParams{Limit: int32(perPageNumber), Offset: int32(perPageNumber * pageNumber)})

	if err != nil {
		http.Error(requestWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	var deviceDtos []DeviceDto
	for _, device := range devices {
		deviceDtos = append(deviceDtos, NewDeviceDto(device))
	}
	if deviceDtos == nil {
		deviceDtos = []DeviceDto{}
	}
	json.NewEncoder(requestWriter).Encode(deviceDtos)

}
