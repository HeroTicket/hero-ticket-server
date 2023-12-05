package main

import (
	"context"
	"fmt"

	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/service/ticket"
	tmongo "github.com/heroticket/internal/service/ticket/repository/mongo"
)

func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	client, err := mongo.New(context.Background(), cfg.MongoUrl)
	if err != nil {
		panic(err)
	}

	db, err := tmongo.New(context.Background(), client, "herotickettest")
	if err != nil {
		panic(err)
	}

	res, err := db.CreateTicketCollection(context.Background(), ticket.CreateTicketCollectionParams{
		ContractAddress: "0x1234567890",
		IssuerAddress:   "0x1234567890",
		Name:            "Test",
		Symbol:          "TST",
		Description:     "Test",
		Organizer:       "Test",
		Location:        "Test",
		Date:            "Test",
		BannerUrl:       "Test",
		TicketUrl:       "Test",
		EthPrice:        "Test",
		TokenPrice:      "Test",
		TotalSupply:     "Test",
		Remaining:       "Test",
		SaleStartAt:     0,
		SaleEndAt:       0,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
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
