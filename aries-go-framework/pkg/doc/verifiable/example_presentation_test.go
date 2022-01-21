/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package verifiable_test

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/jsonld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/suite/ed25519signature2018"
	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util"
	"github.com/hyperledger/aries-framework-go/pkg/doc/util/signature"
	"github.com/hyperledger/aries-framework-go/pkg/doc/verifiable"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
)

// The keys are generated by ed25519.GenerateKey(rand.Reader)
//nolint:gochecknoglobals
var (
	holderPrivKey = ed25519.PrivateKey{10, 192, 72, 230, 66, 255, 51, 97, 14, 57, 149, 164, 232, 251, 31, 164, 168, 82, 239, 155, 253, 223, 111, 148, 165, 76, 60, 17, 3, 63, 76, 192, 61, 133, 23, 17, 77, 132, 169, 196, 47, 203, 19, 71, 145, 144, 92, 145, 131, 101, 36, 251, 89, 216, 117, 140, 132, 226, 78, 187, 59, 58, 200, 255}
	holderPubKey  = ed25519.PublicKey{61, 133, 23, 17, 77, 132, 169, 196, 47, 203, 19, 71, 145, 144, 92, 145, 131, 101, 36, 251, 89, 216, 117, 140, 132, 226, 78, 187, 59, 58, 200, 255}
)

func ExampleParsePresentation() {
	// A Holder sends to the Verifier a verifiable presentation in JWS form.
	vpJWS := "eyJhbGciOiJFZERTQSIsImtpZCI6ImtleS0xIiwidHlwIjoiSldUIn0.eyJpc3MiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJqdGkiOiJ1cm46dXVpZDozOTc4MzQ0Zi04NTk2LTRjM2EtYTk3OC04ZmNhYmEzOTAzYzUiLCJ2cCI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZVByZXNlbnRhdGlvbiIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl0sInZlcmlmaWFibGVDcmVkZW50aWFsIjpbeyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTY2hlbWEiOltdLCJjcmVkZW50aWFsU3ViamVjdCI6eyJkZWdyZWUiOnsidHlwZSI6IkJhY2hlbG9yRGVncmVlIiwidW5pdmVyc2l0eSI6Ik1JVCJ9LCJpZCI6ImRpZDpleGFtcGxlOmViZmViMWY3MTJlYmM2ZjFjMjc2ZTEyZWMyMSIsIm5hbWUiOiJKYXlkZW4gRG9lIiwic3BvdXNlIjoiZGlkOmV4YW1wbGU6YzI3NmUxMmVjMjFlYmZlYjFmNzEyZWJjNmYxIn0sImV4cGlyYXRpb25EYXRlIjoiMjAyMC0wMS0wMVQxOToyMzoyNFoiLCJpZCI6Imh0dHA6Ly9leGFtcGxlLmVkdS9jcmVkZW50aWFscy8xODcyIiwiaXNzdWFuY2VEYXRlIjoiMjAxMC0wMS0wMVQxOToyMzoyNFoiLCJpc3N1ZXIiOnsiaWQiOiJkaWQ6ZXhhbXBsZTo3NmUxMmVjNzEyZWJjNmYxYzIyMWViZmViMWYiLCJuYW1lIjoiRXhhbXBsZSBVbml2ZXJzaXR5In0sInJlZmVyZW5jZU51bWJlciI6OC4zMjk0ODQ3ZSswNywidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19XX19.RlO_1B-7qhQNwo2mmOFUWSa8A6hwaJrtq3q7yJDkKq4k6B-EJ-oyLNM6H_g2_nko2Yg9Im1CiROFm6nK12U_AQ"

	// Holder received and decodes it.
	vp, err := verifiable.ParsePresentation(
		[]byte(vpJWS),
		verifiable.WithPresPublicKeyFetcher(verifiable.SingleKey(holderPubKey, kms.ED25519)),
		verifiable.WithPresJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VP JWS: %w", err))
	}

	// Marshal the VP to JSON to verify the result of decoding.
	vpBytes, err := json.Marshal(vp)
	if err != nil {
		panic(fmt.Errorf("failed to marshal VP to JSON: %w", err))
	}

	fmt.Println(string(vpBytes))

	// Output: {"@context":["https://www.w3.org/2018/credentials/v1","https://www.w3.org/2018/credentials/examples/v1"],"holder":"did:example:ebfeb1f712ebc6f1c276e12ec21","id":"urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c5","type":["VerifiablePresentation","UniversityDegreeCredential"],"verifiableCredential":[{"@context":["https://www.w3.org/2018/credentials/v1","https://www.w3.org/2018/credentials/examples/v1"],"credentialSchema":[],"credentialSubject":{"degree":{"type":"BachelorDegree","university":"MIT"},"id":"did:example:ebfeb1f712ebc6f1c276e12ec21","name":"Jayden Doe","spouse":"did:example:c276e12ec21ebfeb1f712ebc6f1"},"expirationDate":"2020-01-01T19:23:24Z","id":"http://example.edu/credentials/1872","issuanceDate":"2010-01-01T19:23:24Z","issuer":{"id":"did:example:76e12ec712ebc6f1c221ebfeb1f","name":"Example University"},"referenceNumber":83294847,"type":["VerifiableCredential","UniversityDegreeCredential"]}]}
}

func ExamplePresentation_JWTClaims() {
	// The Holder kept the presentation serialized to JSON in her personal verifiable credential wallet.
	vpStrFromWallet := `
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  "id": "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c5",
  "type": [
    "VerifiablePresentation",
    "UniversityDegreeCredential"
  ],
  "verifiableCredential": [
    {
      "@context": [
        "https://www.w3.org/2018/credentials/v1",
        "https://www.w3.org/2018/credentials/examples/v1"
      ],
      "credentialSchema": [],
      "credentialSubject": {
        "degree": {
          "type": "BachelorDegree",
          "university": "MIT"
        },
        "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
        "name": "Jayden Doe",
        "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
      },
      "expirationDate": "2020-01-01T19:23:24Z",
      "id": "http://example.edu/credentials/1872",
      "issuanceDate": "2010-01-01T19:23:24Z",
      "issuer": {
        "id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
        "name": "Example University"
      },
      "referenceNumber": 83294847,
      "type": [
        "VerifiableCredential",
        "UniversityDegreeCredential"
      ]
    }
  ],
  "holder": "did:example:ebfeb1f712ebc6f1c276e12ec21"
}
`

	// The Holder wants to send the presentation to the Verifier in JWS.
	vp, err := verifiable.ParsePresentation([]byte(vpStrFromWallet), verifiable.WithPresDisabledProofCheck(),
		verifiable.WithPresJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VP JSON: %w", err))
	}

	aud := []string{"did:example:4a57546973436f6f6c4a4a57573"}

	jwtClaims, err := vp.JWTClaims(aud, true)
	if err != nil {
		panic(fmt.Errorf("failed to create JWT claims of VP: %w", err))
	}

	signer := signature.GetEd25519Signer(holderPrivKey, holderPubKey)

	jws, err := jwtClaims.MarshalJWS(verifiable.EdDSA, signer, "")
	if err != nil {
		panic(fmt.Errorf("failed to sign VP inside JWT: %w", err))
	}

	fmt.Println(jws)

	// Output: eyJhbGciOiJFZERTQSIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkaWQ6ZXhhbXBsZTo0YTU3NTQ2OTczNDM2ZjZmNmM0YTRhNTc1NzMiLCJpc3MiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJqdGkiOiJ1cm46dXVpZDozOTc4MzQ0Zi04NTk2LTRjM2EtYTk3OC04ZmNhYmEzOTAzYzUiLCJ2cCI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sInR5cGUiOlsiVmVyaWZpYWJsZVByZXNlbnRhdGlvbiIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl0sInZlcmlmaWFibGVDcmVkZW50aWFsIjpbeyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTY2hlbWEiOltdLCJjcmVkZW50aWFsU3ViamVjdCI6eyJkZWdyZWUiOnsidHlwZSI6IkJhY2hlbG9yRGVncmVlIiwidW5pdmVyc2l0eSI6Ik1JVCJ9LCJpZCI6ImRpZDpleGFtcGxlOmViZmViMWY3MTJlYmM2ZjFjMjc2ZTEyZWMyMSIsIm5hbWUiOiJKYXlkZW4gRG9lIiwic3BvdXNlIjoiZGlkOmV4YW1wbGU6YzI3NmUxMmVjMjFlYmZlYjFmNzEyZWJjNmYxIn0sImV4cGlyYXRpb25EYXRlIjoiMjAyMC0wMS0wMVQxOToyMzoyNFoiLCJpZCI6Imh0dHA6Ly9leGFtcGxlLmVkdS9jcmVkZW50aWFscy8xODcyIiwiaXNzdWFuY2VEYXRlIjoiMjAxMC0wMS0wMVQxOToyMzoyNFoiLCJpc3N1ZXIiOnsiaWQiOiJkaWQ6ZXhhbXBsZTo3NmUxMmVjNzEyZWJjNmYxYzIyMWViZmViMWYiLCJuYW1lIjoiRXhhbXBsZSBVbml2ZXJzaXR5In0sInJlZmVyZW5jZU51bWJlciI6ODMyOTQ4NDcsInR5cGUiOlsiVmVyaWZpYWJsZUNyZWRlbnRpYWwiLCJVbml2ZXJzaXR5RGVncmVlQ3JlZGVudGlhbCJdfV19fQ.Rfh63wppLmXFW0nIC1y9aaB9mOJY4D9yQiEcxFVbFjjhjNj6BY9FxzLC67zpl0lwXdIoBeRX_lOhWGOvWo8UAA
}

func ExamplePresentation() {
	// A Holder loads the credential from verifiable credential wallet in order to send to Verifier.
	// She embedded the credential into presentation and sends it to the Verifier in JWS form.
	vcStrFromWallet := `
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  "credentialSubject": {
    "degree": {
      "type": "BachelorDegree",
      "university": "MIT"
    },
    "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
    "name": "Jayden Doe",
    "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
  },
  "expirationDate": "2020-01-01T19:23:24Z",
  "id": "http://example.edu/credentials/1872",
  "issuanceDate": "2010-01-01T19:23:24Z",
  "issuer": {
    "id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
    "name": "Example University"
  },
  "referenceNumber": 83294847,
  "type": [
    "VerifiableCredential",
    "UniversityDegreeCredential"
  ]
}
`

	vc, err := verifiable.ParseCredential([]byte(vcStrFromWallet),
		verifiable.WithJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VC JSON: %w", err))
	}

	vp, err := verifiable.NewPresentation(verifiable.WithCredentials(vc))
	if err != nil {
		panic(fmt.Errorf("failed to build VP from VC: %w", err))
	}

	vp.ID = "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c5"
	vp.Holder = "did:example:ebfeb1f712ebc6f1c276e12ec21"

	aud := []string{"did:example:4a57546973436f6f6c4a4a57573"}

	jwtClaims, err := vp.JWTClaims(aud, true)
	if err != nil {
		panic(fmt.Errorf("failed to create JWT claims of VP: %w", err))
	}

	signer := signature.GetEd25519Signer(holderPrivKey, holderPubKey)

	jws, err := jwtClaims.MarshalJWS(verifiable.EdDSA, signer, "")
	if err != nil {
		panic(fmt.Errorf("failed to sign VP inside JWT: %w", err))
	}

	fmt.Println(jws)

	// Output: eyJhbGciOiJFZERTQSIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkaWQ6ZXhhbXBsZTo0YTU3NTQ2OTczNDM2ZjZmNmM0YTRhNTc1NzMiLCJpc3MiOiJkaWQ6ZXhhbXBsZTplYmZlYjFmNzEyZWJjNmYxYzI3NmUxMmVjMjEiLCJqdGkiOiJ1cm46dXVpZDozOTc4MzQ0Zi04NTk2LTRjM2EtYTk3OC04ZmNhYmEzOTAzYzUiLCJ2cCI6eyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSJdLCJ0eXBlIjoiVmVyaWZpYWJsZVByZXNlbnRhdGlvbiIsInZlcmlmaWFibGVDcmVkZW50aWFsIjpbeyJAY29udGV4dCI6WyJodHRwczovL3d3dy53My5vcmcvMjAxOC9jcmVkZW50aWFscy92MSIsImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL2V4YW1wbGVzL3YxIl0sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJ0eXBlIjoiQmFjaGVsb3JEZWdyZWUiLCJ1bml2ZXJzaXR5IjoiTUlUIn0sImlkIjoiZGlkOmV4YW1wbGU6ZWJmZWIxZjcxMmViYzZmMWMyNzZlMTJlYzIxIiwibmFtZSI6IkpheWRlbiBEb2UiLCJzcG91c2UiOiJkaWQ6ZXhhbXBsZTpjMjc2ZTEyZWMyMWViZmViMWY3MTJlYmM2ZjEifSwiZXhwaXJhdGlvbkRhdGUiOiIyMDIwLTAxLTAxVDE5OjIzOjI0WiIsImlkIjoiaHR0cDovL2V4YW1wbGUuZWR1L2NyZWRlbnRpYWxzLzE4NzIiLCJpc3N1YW5jZURhdGUiOiIyMDEwLTAxLTAxVDE5OjIzOjI0WiIsImlzc3VlciI6eyJpZCI6ImRpZDpleGFtcGxlOjc2ZTEyZWM3MTJlYmM2ZjFjMjIxZWJmZWIxZiIsIm5hbWUiOiJFeGFtcGxlIFVuaXZlcnNpdHkifSwicmVmZXJlbmNlTnVtYmVyIjo4MzI5NDg0NywidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19XX19.DMayxVTjX-tKwemmIuoJvxw8A0Oj5KMM1xgKF_SaFO4GQHAspQEDT70RJrW37WDHaYnFyVAimTLlGkaxKic-Dg
}

func ExamplePresentation_two() {
	// Holder wants to send 2 credentials to Verifier.
	// The first VC is created on fly (or just decoded using ParseCredential).
	vc := &verifiable.Credential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		ID: "http://example.edu/credentials/1872",
		Types: []string{
			"VerifiableCredential",
			"UniversityDegreeCredential",
		},
		Subject: UniversityDegreeSubject{
			ID:     "did:example:ebfeb1f712ebc6f1c276e12ec21",
			Name:   "Jayden Doe",
			Spouse: "did:example:c276e12ec21ebfeb1f712ebc6f1",
			Degree: UniversityDegree{
				Type:       "BachelorDegree",
				University: "MIT",
			},
		},
		Issuer: verifiable.Issuer{
			ID:           "did:example:76e12ec712ebc6f1c221ebfeb1f",
			CustomFields: verifiable.CustomFields{"name": "Example University"},
		},
		Issued:  util.NewTime(issued),
		Expired: util.NewTime(expired),
		Schemas: []verifiable.TypedID{},
		CustomFields: map[string]interface{}{
			"referenceNumber": 83294847,
		},
	}

	vcStr := `
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  "id": "http://example.edu/credentials/58473",
  "type": ["VerifiableCredential", "AlumniCredential"],
  "issuer": "https://example.edu/issuers/14",
  "issuanceDate": "2010-01-01T19:23:24Z",
  "credentialSubject": {
    "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
    "alumniOf": "Example University"
  },
  "proof": {
    "type": "RsaSignature2018"
  }
}
`

	// The second VC is provided in JWS form (e.g. kept in the wallet in that form).
	vcJWS := "eyJhbGciOiJFZERTQSIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc5MDY2MDQsImlhdCI6MTI2MjM3MzgwNCwiaXNzIjoiZGlkOmV4YW1wbGU6NzZlMTJlYzcxMmViYzZmMWMyMjFlYmZlYjFmIiwianRpIjoiaHR0cDovL2V4YW1wbGUuZWR1L2NyZWRlbnRpYWxzLzE4NzIiLCJuYmYiOjEyNjIzNzM4MDQsInN1YiI6ImRpZDpleGFtcGxlOmViZmViMWY3MTJlYmM2ZjFjMjc2ZTEyZWMyMSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIiwiaHR0cHM6Ly93d3cudzMub3JnLzIwMTgvY3JlZGVudGlhbHMvZXhhbXBsZXMvdjEiXSwiY3JlZGVudGlhbFNjaGVtYSI6W10sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJ0eXBlIjoiQmFjaGVsb3JEZWdyZWUiLCJ1bml2ZXJzaXR5IjoiTUlUIn0sImlkIjoiZGlkOmV4YW1wbGU6ZWJmZWIxZjcxMmViYzZmMWMyNzZlMTJlYzIxIiwibmFtZSI6IkpheWRlbiBEb2UiLCJzcG91c2UiOiJkaWQ6ZXhhbXBsZTpjMjc2ZTEyZWMyMWViZmViMWY3MTJlYmM2ZjEifSwiaXNzdWVyIjp7Im5hbWUiOiJFeGFtcGxlIFVuaXZlcnNpdHkifSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19fQ.AHn2A2q5DL1heX3_izq_2yrsBDhoZ6BGGKhoRvhfMnMUuuOnBOdekdTg-dfUMJgipXRql_6WzBUIj4wTFehXCw"

	vc2, err := verifiable.ParseCredential([]byte(vcStr),
		verifiable.WithJSONLDDocumentLoader(getJSONLDDocumentLoader()),
		verifiable.WithDisabledProofCheck())
	if err != nil {
		panic(fmt.Errorf("failed to decode VC JSON: %w", err))
	}

	vp, err := verifiable.NewPresentation(verifiable.WithCredentials(vc),
		verifiable.WithJWTCredentials(vcJWS), verifiable.WithCredentials(vc2))
	if err != nil {
		panic(fmt.Errorf("failed to set credentials of VP: %w", err))
	}

	vp.ID = "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c"
	vp.Holder = "did:example:ebfeb1f712ebc6f1c276e12ec21"

	vpBytes, err := json.MarshalIndent(vp, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Print(string(vpBytes))

	// Output:
	// {
	// 	"@context": [
	// 		"https://www.w3.org/2018/credentials/v1"
	// 	],
	// 	"holder": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	// 	"id": "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c",
	// 	"type": "VerifiablePresentation",
	// 	"verifiableCredential": [
	// 		{
	// 			"@context": [
	// 				"https://www.w3.org/2018/credentials/v1",
	// 				"https://www.w3.org/2018/credentials/examples/v1"
	// 			],
	// 			"credentialSubject": {
	// 				"degree": {
	// 					"type": "BachelorDegree",
	// 					"university": "MIT"
	// 				},
	// 				"id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	// 				"name": "Jayden Doe",
	// 				"spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
	// 			},
	// 			"expirationDate": "2020-01-01T19:23:24Z",
	// 			"id": "http://example.edu/credentials/1872",
	// 			"issuanceDate": "2010-01-01T19:23:24Z",
	// 			"issuer": {
	// 				"id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
	// 				"name": "Example University"
	// 			},
	// 			"referenceNumber": 83294847,
	// 			"type": [
	// 				"VerifiableCredential",
	// 				"UniversityDegreeCredential"
	// 			]
	// 		},
	// 		"eyJhbGciOiJFZERTQSIsImtpZCI6IiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc5MDY2MDQsImlhdCI6MTI2MjM3MzgwNCwiaXNzIjoiZGlkOmV4YW1wbGU6NzZlMTJlYzcxMmViYzZmMWMyMjFlYmZlYjFmIiwianRpIjoiaHR0cDovL2V4YW1wbGUuZWR1L2NyZWRlbnRpYWxzLzE4NzIiLCJuYmYiOjEyNjIzNzM4MDQsInN1YiI6ImRpZDpleGFtcGxlOmViZmViMWY3MTJlYmM2ZjFjMjc2ZTEyZWMyMSIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIiwiaHR0cHM6Ly93d3cudzMub3JnLzIwMTgvY3JlZGVudGlhbHMvZXhhbXBsZXMvdjEiXSwiY3JlZGVudGlhbFNjaGVtYSI6W10sImNyZWRlbnRpYWxTdWJqZWN0Ijp7ImRlZ3JlZSI6eyJ0eXBlIjoiQmFjaGVsb3JEZWdyZWUiLCJ1bml2ZXJzaXR5IjoiTUlUIn0sImlkIjoiZGlkOmV4YW1wbGU6ZWJmZWIxZjcxMmViYzZmMWMyNzZlMTJlYzIxIiwibmFtZSI6IkpheWRlbiBEb2UiLCJzcG91c2UiOiJkaWQ6ZXhhbXBsZTpjMjc2ZTEyZWMyMWViZmViMWY3MTJlYmM2ZjEifSwiaXNzdWVyIjp7Im5hbWUiOiJFeGFtcGxlIFVuaXZlcnNpdHkifSwidHlwZSI6WyJWZXJpZmlhYmxlQ3JlZGVudGlhbCIsIlVuaXZlcnNpdHlEZWdyZWVDcmVkZW50aWFsIl19fQ.AHn2A2q5DL1heX3_izq_2yrsBDhoZ6BGGKhoRvhfMnMUuuOnBOdekdTg-dfUMJgipXRql_6WzBUIj4wTFehXCw",
	// 		{
	// 			"@context": [
	// 				"https://www.w3.org/2018/credentials/v1",
	// 				"https://www.w3.org/2018/credentials/examples/v1"
	// 			],
	// 			"credentialSubject": {
	// 				"alumniOf": "Example University",
	// 				"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"
	// 			},
	// 			"id": "http://example.edu/credentials/58473",
	// 			"issuanceDate": "2010-01-01T19:23:24Z",
	// 			"issuer": "https://example.edu/issuers/14",
	// 			"proof": {
	// 				"type": "RsaSignature2018"
	// 			},
	// 			"type": [
	// 				"VerifiableCredential",
	// 				"AlumniCredential"
	// 			]
	// 		}
	// 	]
	// }
}

func ExamplePresentation_MarshalJSON() {
	vc := &verifiable.Credential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		ID: "http://example.edu/credentials/1872",
		Types: []string{
			"VerifiableCredential",
			"UniversityDegreeCredential",
		},
		Subject: UniversityDegreeSubject{
			ID:     "did:example:ebfeb1f712ebc6f1c276e12ec21",
			Name:   "Jayden Doe",
			Spouse: "did:example:c276e12ec21ebfeb1f712ebc6f1",
			Degree: UniversityDegree{
				Type:       "BachelorDegree",
				University: "MIT",
			},
		},
		Issuer: verifiable.Issuer{
			ID:           "did:example:76e12ec712ebc6f1c221ebfeb1f",
			CustomFields: verifiable.CustomFields{"name": "Example University"},
		},
		Issued:  util.NewTime(issued),
		Expired: util.NewTime(expired),
		Schemas: []verifiable.TypedID{},
		CustomFields: map[string]interface{}{
			"referenceNumber": 83294847,
		},
	}

	vp, err := verifiable.NewPresentation(verifiable.WithCredentials(vc))
	if err != nil {
		panic(fmt.Errorf("failed to set credentials of VP: %w", err))
	}

	vp.ID = "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c"
	vp.Holder = "did:example:ebfeb1f712ebc6f1c276e12ec21"

	// json.MarshalIndent() calls Presentation.MarshalJSON()
	vpJSON, err := json.MarshalIndent(vp, "", "\t")
	if err != nil {
		panic(fmt.Errorf("failed to marshal VP to JSON: %w", err))
	}

	fmt.Println(string(vpJSON))

	// Output:
	// {
	// 	"@context": [
	// 		"https://www.w3.org/2018/credentials/v1"
	// 	],
	// 	"holder": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	// 	"id": "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c",
	// 	"type": "VerifiablePresentation",
	// 	"verifiableCredential": [
	// 		{
	// 			"@context": [
	// 				"https://www.w3.org/2018/credentials/v1",
	// 				"https://www.w3.org/2018/credentials/examples/v1"
	// 			],
	// 			"credentialSubject": {
	// 				"degree": {
	// 					"type": "BachelorDegree",
	// 					"university": "MIT"
	// 				},
	// 				"id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	// 				"name": "Jayden Doe",
	// 				"spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
	// 			},
	// 			"expirationDate": "2020-01-01T19:23:24Z",
	// 			"id": "http://example.edu/credentials/1872",
	// 			"issuanceDate": "2010-01-01T19:23:24Z",
	// 			"issuer": {
	// 				"id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
	// 				"name": "Example University"
	// 			},
	// 			"referenceNumber": 83294847,
	// 			"type": [
	// 				"VerifiableCredential",
	// 				"UniversityDegreeCredential"
	// 			]
	// 		}
	// 	]
	// }
}

//nolint:gocyclo
func ExamplePresentation_MarshalledCredentials() {
	vc := verifiable.Credential{
		Context: []string{
			"https://www.w3.org/2018/credentials/v1",
			"https://www.w3.org/2018/credentials/examples/v1",
		},
		ID: "http://example.edu/credentials/1872",
		Types: []string{
			"VerifiableCredential",
			"UniversityDegreeCredential",
		},
		Subject: UniversityDegreeSubject{
			ID:     "did:example:ebfeb1f712ebc6f1c276e12ec21",
			Name:   "Jayden Doe",
			Spouse: "did:example:c276e12ec21ebfeb1f712ebc6f1",
			Degree: UniversityDegree{
				Type:       "BachelorDegree",
				University: "MIT",
			},
		},
		Issuer: verifiable.Issuer{
			ID:           "did:example:76e12ec712ebc6f1c221ebfeb1f",
			CustomFields: verifiable.CustomFields{"name": "Example University"},
		},
		Issued:  util.NewTime(issued),
		Expired: util.NewTime(expired),
		Schemas: []verifiable.TypedID{},
	}

	// Put JWS form of VC into VP.
	vcJWTClaims, err := vc.JWTClaims(true)
	if err != nil {
		panic(fmt.Errorf("failed to set credentials of VP: %w", err))
	}

	issuerSigner := signature.GetEd25519Signer(issuerPrivKey, issuerPubKey)

	vcJWS, err := vcJWTClaims.MarshalJWS(verifiable.EdDSA, issuerSigner, "i-kid")
	if err != nil {
		panic(fmt.Errorf("failed to sign VC JWT: %w", err))
	}

	vp, err := verifiable.NewPresentation(verifiable.WithJWTCredentials(vcJWS))
	if err != nil {
		panic(fmt.Errorf("failed to set credentials of VP: %w", err))
	}

	vp.ID = "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c"
	vp.Holder = "did:example:ebfeb1f712ebc6f1c276e12ec21"

	// Marshal VP to JWS as well.

	vpJWTClaims, err := vp.JWTClaims(nil, true)
	if err != nil {
		panic(fmt.Errorf("failed to create JWT claims of VP: %w", err))
	}

	holderSigner := signature.GetEd25519Signer(holderPrivKey, holderPubKey)

	vpJWS, err := vpJWTClaims.MarshalJWS(verifiable.EdDSA, holderSigner, "h-kid")
	if err != nil {
		panic(fmt.Errorf("failed to sign VP inside JWT: %w", err))
	}

	// Decode VP from JWS.
	// Note that VC-s inside will be decoded as well. If they are JWS, their signature is verified
	// and thus we need to make sure the public key fetcher can retrieve the public key.
	vp, err = verifiable.ParsePresentation(
		[]byte(vpJWS),
		verifiable.WithPresPublicKeyFetcher(func(issuerID, keyID string) (*verifier.PublicKey, error) {
			switch issuerID {
			case "did:example:76e12ec712ebc6f1c221ebfeb1f":
				return &verifier.PublicKey{
					Type:  kms.ED25519,
					Value: issuerPubKey,
				}, nil
			case "did:example:ebfeb1f712ebc6f1c276e12ec21":
				return &verifier.PublicKey{
					Type:  kms.ED25519,
					Value: holderPubKey,
				}, nil
			default:
				return nil, fmt.Errorf("unexpected key: %s", keyID)
			}
		}), verifiable.WithPresJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VP JWS: %w", err))
	}

	// Get credentials in binary form.
	vpCreds, err := vp.MarshalledCredentials()
	if err != nil {
		panic(fmt.Errorf("failed to get marshalled credentials from decoded presentation: %w", err))
	}

	if len(vpCreds) != 1 {
		panic("Expected 1 credential inside presentation")
	}

	// Decoded credential. Note that no public key fetcher is passed as the VC was already decoded (and proof verified)
	// when VP was decoded.
	vcDecoded, err := verifiable.ParseCredential(vpCreds[0],
		verifiable.WithJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VC: %w", err))
	}

	vcDecodedJSON, err := json.MarshalIndent(vcDecoded, "", "\t")
	if err != nil {
		panic(fmt.Errorf("failed to marshal VP to JSON: %w", err))
	}

	fmt.Println(string(vcDecodedJSON))

	// Output: {
	//	"@context": [
	//		"https://www.w3.org/2018/credentials/v1",
	//		"https://www.w3.org/2018/credentials/examples/v1"
	//	],
	//	"credentialSubject": {
	//		"degree": {
	//			"type": "BachelorDegree",
	//			"university": "MIT"
	//		},
	//		"id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	//		"name": "Jayden Doe",
	//		"spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
	//	},
	//	"expirationDate": "2020-01-01T19:23:24Z",
	//	"id": "http://example.edu/credentials/1872",
	//	"issuanceDate": "2010-01-01T19:23:24Z",
	//	"issuer": {
	//		"id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
	//		"name": "Example University"
	//	},
	//	"type": [
	//		"VerifiableCredential",
	//		"UniversityDegreeCredential"
	//	]
	//}
}

//nolint:gocyclo
func ExamplePresentation_AddLinkedDataProof() {
	// 1. ISSUER issues a VC.
	vcToIssue := `
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://www.w3.org/2018/credentials/examples/v1"
  ],
  "credentialSubject": {
    "degree": {
      "type": "BachelorDegree",
      "university": "MIT"
    },
    "id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
    "name": "Jayden Doe",
    "spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
  },
  "expirationDate": "2020-01-01T19:23:24Z",
  "id": "http://example.edu/credentials/1872",
  "issuanceDate": "2010-01-01T19:23:24Z",
  "issuer": {
    "id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
    "name": "Example University"
  },
  "type": [
    "VerifiableCredential",
    "UniversityDegreeCredential"
  ]
}
`

	issuedVC, err := verifiable.ParseCredential([]byte(vcToIssue),
		verifiable.WithJSONLDDocumentLoader(getJSONLDDocumentLoader()),
		verifiable.WithDisabledProofCheck())
	if err != nil {
		panic(fmt.Errorf("failed to decode VC JSON: %w", err))
	}

	issuerSigner := signature.GetEd25519Signer(issuerPrivKey, issuerPubKey)

	err = issuedVC.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
		Created:                 &issued,
		SignatureType:           "Ed25519Signature2018",
		Suite:                   ed25519signature2018.New(suite.WithSigner(issuerSigner)),
		SignatureRepresentation: verifiable.SignatureJWS,
		VerificationMethod:      "did:example:123456#key1",
	}, jsonld.WithDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to add linked data proof: %w", err))
	}

	issuedVCBytes, err := issuedVC.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("failed to marshal VC to JSON: %w", err))
	}

	// 2. ISSUER creates a VP with the VC enclosed.
	vcFromHolderWallet, err := verifiable.ParseCredential(issuedVCBytes,
		verifiable.WithJSONLDDocumentLoader(getJSONLDDocumentLoader()),
		verifiable.WithDisabledProofCheck())
	if err != nil {
		panic(fmt.Errorf("failed to decode VC JSON: %w", err))
	}

	vpToVerify, err := verifiable.NewPresentation(verifiable.WithCredentials(vcFromHolderWallet))
	if err != nil {
		panic(fmt.Errorf("failed to build VP from VC: %w", err))
	}

	vpToVerify.Holder = "did:example:ebfeb1f712ebc6f1c276e12ec22"
	vpToVerify.ID = "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c6"

	holderVerifier := signature.GetEd25519Signer(holderPrivKey, holderPubKey)

	err = vpToVerify.AddLinkedDataProof(&verifiable.LinkedDataProofContext{
		Created:                 &issued,
		SignatureType:           "Ed25519Signature2018",
		Suite:                   ed25519signature2018.New(suite.WithSigner(holderVerifier)),
		SignatureRepresentation: verifiable.SignatureJWS,
		VerificationMethod:      "did:example:987654#key1",
	}, jsonld.WithDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to add linked data proof: %w", err))
	}

	vpJSONWithProof, err := vpToVerify.MarshalJSON()
	if err != nil {
		panic(fmt.Errorf("failed to marshal VP to JSON: %w", err))
	}

	// 3. VERIFIER verifies the presentation.
	ed25519Suite := ed25519signature2018.New(suite.WithVerifier(ed25519signature2018.NewPublicKeyVerifier()))

	vp, err := verifiable.ParsePresentation(vpJSONWithProof,
		verifiable.WithPresPublicKeyFetcher(func(issuerID, keyID string) (*verifier.PublicKey, error) {
			// both VP and enclosed VC signatures are verified, so we need to provide key resolving for all
			switch issuerID {
			case "did:example:123456": // issuer
				return &verifier.PublicKey{
					Type:  "Ed25519Signature2018",
					Value: issuerPubKey,
				}, nil

			case "did:example:987654":
				return &verifier.PublicKey{
					Type:  "Ed25519Signature2018",
					Value: holderPubKey,
				}, nil
			}

			return nil, errors.New("unsupported issuer")
		}),
		verifiable.WithPresEmbeddedSignatureSuites(ed25519Suite),
		verifiable.WithPresJSONLDDocumentLoader(getJSONLDDocumentLoader()))
	if err != nil {
		panic(fmt.Errorf("failed to decode VP JWS: %w", err))
	}

	vpJSON, err := json.MarshalIndent(vp, "", "\t")
	if err != nil {
		panic(fmt.Errorf("failed to marshal VC to JSON: %w", err))
	}

	fmt.Println(string(vpJSON))

	// Output: {
	//	"@context": [
	//		"https://www.w3.org/2018/credentials/v1"
	//	],
	//	"holder": "did:example:ebfeb1f712ebc6f1c276e12ec22",
	//	"id": "urn:uuid:3978344f-8596-4c3a-a978-8fcaba3903c6",
	//	"proof": {
	//		"created": "2010-01-01T19:23:24Z",
	//		"jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..8stDRasAcYjkQiqiczyFJdkff8VJIF3Lbaq5BNTaC-PcvJHGo2Xja8GTsHByTOx7QNCwC3bNiboPgfXtmm8aBA",
	//		"proofPurpose": "assertionMethod",
	//		"type": "Ed25519Signature2018",
	//		"verificationMethod": "did:example:987654#key1"
	//	},
	//	"type": "VerifiablePresentation",
	//	"verifiableCredential": [
	//		{
	//			"@context": [
	//				"https://www.w3.org/2018/credentials/v1",
	//				"https://www.w3.org/2018/credentials/examples/v1"
	//			],
	//			"credentialSubject": {
	//				"degree": {
	//					"type": "BachelorDegree",
	//					"university": "MIT"
	//				},
	//				"id": "did:example:ebfeb1f712ebc6f1c276e12ec21",
	//				"name": "Jayden Doe",
	//				"spouse": "did:example:c276e12ec21ebfeb1f712ebc6f1"
	//			},
	//			"expirationDate": "2020-01-01T19:23:24Z",
	//			"id": "http://example.edu/credentials/1872",
	//			"issuanceDate": "2010-01-01T19:23:24Z",
	//			"issuer": {
	//				"id": "did:example:76e12ec712ebc6f1c221ebfeb1f",
	//				"name": "Example University"
	//			},
	//			"proof": {
	//				"created": "2010-01-01T19:23:24Z",
	//				"jws": "eyJhbGciOiJFZERTQSIsImI2NCI6ZmFsc2UsImNyaXQiOlsiYjY0Il19..mQCxgQDvAYI-2YYCkHHe-at9eNI_wN03R6CRyjycb3CnfPWezbo6zEGe94W2AdYsBhC_Zzedcn_ZKgccMYFnCQ",
	//				"proofPurpose": "assertionMethod",
	//				"type": "Ed25519Signature2018",
	//				"verificationMethod": "did:example:123456#key1"
	//			},
	//			"type": [
	//				"VerifiableCredential",
	//				"UniversityDegreeCredential"
	//			]
	//		}
	//	]
	//}
}
