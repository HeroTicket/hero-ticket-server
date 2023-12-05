package main

import (
	"context"
	"net/http"
	"os"

	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/service/ipfs"
)

func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	svc := ipfs.New(ipfs.IpfsServiceConfig{
		ApiKey: cfg.Ipfs.ApiKey,
		Secret: cfg.Ipfs.Secret,
		Client: http.DefaultClient,
	})

	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	resp, err := svc.PinFile(context.Background(), f, "엄준식")
	if err != nil {
		panic(err)
	}

	println(resp.IpfsHash)
}

/*
func main() {
	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r, w := io.Pipe()

	fmt.Println("Pipe Created")

	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		part, err := m.CreateFormFile("file", f.Name())
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(part, f); err != nil {
			panic(err)
		}
	}()

	fmt.Println("Create Request")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, "https://api.pinata.cloud/pinning/pinFileToIPFS", r)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", m.FormDataContentType())

	fmt.Println("Content-Type:", m.FormDataContentType())

	req.Header.Add("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySW5mb3JtYXRpb24iOnsiaWQiOiIwY2Y0MTZlMi1kZWJiLTQ3MjgtYmUwOC0zMWEwMWNhNTAwOTIiLCJlbWFpbCI6ImNyZXdlMTc0NkBoYW55YW5nLmFjLmtyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInBpbl9wb2xpY3kiOnsicmVnaW9ucyI6W3siaWQiOiJGUkExIiwiZGVzaXJlZFJlcGxpY2F0aW9uQ291bnQiOjF9LHsiaWQiOiJOWUMxIiwiZGVzaXJlZFJlcGxpY2F0aW9uQ291bnQiOjF9XSwidmVyc2lvbiI6MX0sIm1mYV9lbmFibGVkIjpmYWxzZSwic3RhdHVzIjoiQUNUSVZFIn0sImF1dGhlbnRpY2F0aW9uVHlwZSI6InNjb3BlZEtleSIsInNjb3BlZEtleUtleSI6IjAxN2U0ZTg2ODY1YmRhMDNkMDFiIiwic2NvcGVkS2V5U2VjcmV0IjoiNDM3N2RmMWU0MWFlZDBiNGM1OWFhMDQxZWVjODEwNDExNGRlMzllMmRjNzRmMTA1YjkxZWY3MzQ4OGEzMDM4OSIsImlhdCI6MTcwMTI0NzUzOH0.1hopoeAPpdDUCeFp-TbkWAbnhgADVYZ6KVw9RG7tQHQ")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Status Code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("Error:", string(b))
		return
	}

	var data ipfs.PinFileResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	fmt.Println(data)
}
*/
/*
func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	f, err := os.Open("test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pnt := pinata.Pinata{
		Apikey: cfg.Ipfs.ApiKey,
		Secret: cfg.Ipfs.Secret,
	}

	hash, err := pnt.PinWithReader(f)
	if err != nil {
		panic(err)
	}

	println(hash)
}
*/
