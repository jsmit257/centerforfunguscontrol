package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/jsmit257/huautla/types"

	"github.com/stretchr/testify/require"
)

var (
	strains     []types.Strain
	sampledata  bytes.Buffer
	contentType string
)

func init() {
	const samplefile = "../../tests/data/sample.png"

	w := multipart.NewWriter(&sampledata)

	if sample, err := os.Open(samplefile); err != nil {
		panic(err)
	} else if sample == nil {
		panic(fmt.Errorf("sample image data is nil"))
	} else if fw, err := w.CreateFormFile("file", samplefile); err != nil {
		panic(err)
	} else if _, err = io.Copy(fw, sample); err != nil {
		panic(err)
	}

	w.Close()

	contentType = w.FormDataContentType()
}

func Test_HappyStrain(t *testing.T) {
	url := fmt.Sprintf(`http://%s:%d/strain`, cfg.HTTPHost, cfg.HTTPPort)

	for _, s := range []types.Strain{
		{Name: "Morel", Species: "M.anatolica", Vendor: vendors[0]},
		{Name: "Hens o' the Wood", Species: "G.frondosa", Vendor: vendors[2]},
		{Name: "Reishi", Species: "G.lingzhi", Vendor: vendors[1]},
		{Name: "Turkey Tail", Species: "T.versicolor", Vendor: vendors[3]},
		{Name: "Scooby Snacks", Species: "Vanilla Wafers", Vendor: vendors[1]},
		{Name: "Shitake", Species: "L.edodes", Vendor: vendors[0]},
		{Name: "Chestnut", Species: "P.adiposa", Vendor: vendors[2]},
		{Name: "Hericium", Species: "H.abietis", Vendor: vendors[0]},
	} {
		b, err := json.Marshal(s)
		require.Nil(t, err)

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
		require.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		require.Nil(t, err)
		require.Equal(t, http.StatusCreated, res.StatusCode)

		b, err = io.ReadAll(res.Body)
		require.Nil(t, err)

		err = json.Unmarshal(b, &s)
		require.Nil(t, err)

		strains = append(strains, s)
	}
}

func createPhoto(t *testing.T, id types.UUID) types.UUID {
	url := fmt.Sprintf(`http://%s:%d/photos`, cfg.HTTPHost, cfg.HTTPPort)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", url, id),
		bytes.NewReader(sampledata.Bytes()))
	require.Nil(t, err)

	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode, "strain: %v", id)

	var b []byte
	b, err = io.ReadAll(res.Body)
	require.Nil(t, err)

	var photos []types.Photo
	err = json.Unmarshal(b, &photos)
	require.Nil(t, err)
	require.Equal(t, 1, len(photos))

	return photos[0].UUID
}

func createNote(t *testing.T, id types.UUID, note string) {
	url := fmt.Sprintf(`http://%s:%d/notes`, cfg.HTTPHost, cfg.HTTPPort)

	b, err := json.Marshal(types.Note{Note: note})
	require.Nil(t, err)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/%s", url, id),
		bytes.NewReader(b))
	require.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	b, err = io.ReadAll(res.Body)
	require.Nil(t, err)

	var notes []types.Note
	err = json.Unmarshal(b, &notes)
	require.Nil(t, err)
	require.Equal(t, 1, len(notes))
}

func Test_HappyStrainPhoto(t *testing.T) {
	for _, s := range strains {
		createNote(t, createPhoto(t, s.UUID), s.Name)
	}
}
