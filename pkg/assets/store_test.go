// Copyright 2020 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assets

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

var (
	caPEM = `-----BEGIN CERTIFICATE-----
MIIB4zCCAY2gAwIBAgIUf+9T+SQuY7RzRfLrT/m3ZLZa/nswDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAgFw0yMDEwMTkxMzA1MDlaGA8yMTIw
MDkyNTEzMDUwOVowRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUx
ITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDBcMA0GCSqGSIb3DQEB
AQUAA0sAMEgCQQDbXwmz6fkHnfs3p5dirgW/m5G1eOSddS8atIwhOzaYSNG03/Z4
P6HWCGDCgUg77fOsX+tzYWkXy0T+GwQrTLDdAgMBAAGjUzBRMB0GA1UdDgQWBBTC
CNvaPTFE1Xt5WUREDoF/mTOg7DAfBgNVHSMEGDAWgBTCCNvaPTFE1Xt5WUREDoF/
mTOg7DAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA0EAzhzA2n5nSnka
k9iw9ZHayRBSgnGAYKFdiGyvceKPzR3LJ8vMdGeYh/TSHHgZ4QSam/J7vHWCkJmc
7c98vpkIaw==
-----END CERTIFICATE-----`

	certPEM = `-----BEGIN CERTIFICATE-----
MIIBiTCCATMCFCgn66sq14Tsx6iP8nRdP4/uiguXMA0GCSqGSIb3DQEBCwUAMEUx
CzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRl
cm5ldCBXaWRnaXRzIFB0eSBMdGQwIBcNMjAxMDE5MTMwNTI5WhgPMjEyMDA5MjUx
MzA1MjlaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYD
VQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEFAANL
ADBIAkEA1wNGN6zrF4eBpW/LcNt3Qxy9bZZss6c/pUy5V4n2O+tZZuvKXF3Q6g4+
fOZ5xgqzqPgg2UzrG1Mmt/Ol4UikZQIDAQABMA0GCSqGSIb3DQEBCwUAA0EAGsWD
5UlmIIbFOi50jqNE3KitIwbPuY8nYR8pS2HYSE+eVKpGFmmzIRXkb4ZmdVymI+vG
B9nfCt+guZqCLxZMDQ==
-----END CERTIFICATE-----`

	keyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIBPAIBAAJBANcDRjes6xeHgaVvy3Dbd0McvW2WbLOnP6VMuVeJ9jvrWWbrylxd
0OoOPnzmecYKs6j4INlM6xtTJrfzpeFIpGUCAwEAAQJAMhPxJsZ/ett0trNzDrYO
8PKgrAV9C9rIWBemk1zunMWmmtBt295sEK555iedWanANhTYKlaezUXMBZaoHIhc
AQIhAPB6QM5fGEsH1VSXEgaSb/EewQLFGjkWj9DtFtwOtmWpAiEA5OQ7NTVq9ULq
6qAI/JJ6qVGCjS/bmUQD2aBrUUhdxl0CIQDrOvsno/fUdS4ll70nNplPqICu3/Ud
wMcfXLwOuEmNOQIhAMSYi4o+IWobWe7AGjfmEFkR25ItAu73jl8D/GlKQNE5AiEA
hvBlhCknnq89u57O41ID6Mqxz3bRxNxpkqhfMyVWcVU=
-----END RSA PRIVATE KEY-----`
)

func TestGetSecretKey(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
	)

	for _, tc := range []struct {
		ns           string
		selectedName string
		selectedKey  string

		err      bool
		expected string
	}{
		{
			ns:           "ns1",
			selectedName: "secret",
			selectedKey:  "key1",

			expected: "val1",
		},
		// Wrong namespace.
		{
			ns:           "ns2",
			selectedName: "secret",
			selectedKey:  "key1",

			err: true,
		},
		// Wrong name.
		{
			ns:           "ns1",
			selectedName: "secreet",
			selectedKey:  "key1",

			err: true,
		},
		// Wrong key.
		{
			ns:           "ns1",
			selectedName: "secret",
			selectedKey:  "key2",

			err: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			sel := v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: tc.selectedName,
				},
				Key: tc.selectedKey,
			}

			s, err := store.GetSecretKey(context.Background(), tc.ns, sel)

			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.expected, s, "expecting %q, got %q", tc.expected, s)
		})
	}
}

func TestAddBasicAuth(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
				"key2": []byte("val2"),
			},
		},
	)

	for _, tc := range []struct {
		ns                   string
		selectedUserName     string
		selectedUserKey      string
		selectedPasswordName string
		selectedPasswordKey  string

		err              bool
		expectedUser     string
		expectedPassword string
	}{
		{
			ns:                   "ns1",
			selectedUserName:     "secret",
			selectedUserKey:      "key1",
			selectedPasswordName: "secret",
			selectedPasswordKey:  "key2",

			expectedUser:     "val1",
			expectedPassword: "val2",
		},
		// Wrong namespace.
		{
			ns:                   "ns2",
			selectedUserName:     "secret",
			selectedUserKey:      "key1",
			selectedPasswordName: "secret",
			selectedPasswordKey:  "key2",

			err: true,
		},
		// Wrong name for username selector.
		{
			ns:                   "ns1",
			selectedUserName:     "secreet",
			selectedUserKey:      "key1",
			selectedPasswordName: "secret",
			selectedPasswordKey:  "key2",

			err: true,
		},
		// Wrong key for username selector.
		{
			ns:                   "ns1",
			selectedUserName:     "secret",
			selectedUserKey:      "key3",
			selectedPasswordName: "secret",
			selectedPasswordKey:  "key2",

			err: true,
		},
		// Wrong name for password selector.
		{
			ns:                   "ns1",
			selectedUserName:     "secret",
			selectedUserKey:      "key1",
			selectedPasswordName: "secreet",
			selectedPasswordKey:  "key2",

			err: true,
		},
		// Wrong key for password selector.
		{
			ns:                   "ns1",
			selectedUserName:     "secret",
			selectedUserKey:      "key1",
			selectedPasswordName: "secret",
			selectedPasswordKey:  "key3",

			err: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			basicAuth := &monitoringv1.BasicAuth{
				Username: v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: tc.selectedUserName,
					},
					Key: tc.selectedUserKey,
				},
				Password: v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: tc.selectedPasswordName,
					},
					Key: tc.selectedPasswordKey,
				},
			}

			err := store.AddBasicAuth(context.Background(), tc.ns, basicAuth)

			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			b, err := store.ForNamespace(tc.ns).GetSecretKey(basicAuth.Password)
			require.NoError(t, err)

			require.Equal(t, tc.expectedPassword, string(b), "expecting password value %q, got %q", tc.expectedPassword, string(b))

			b, err = store.ForNamespace(tc.ns).GetSecretKey(basicAuth.Username)
			require.NoError(t, err)

			require.Equal(t, tc.expectedUser, string(b), "expecting username value %q, got %q", tc.expectedUser, string(b))
		})
	}
}

func TestProxyCongfig(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"proxyA": []byte("proxyA"),
				"proxyB": []byte("proxyB"),
				"proxyC": []byte("proxyC"),
			},
		},
	)

	for _, tc := range []struct {
		ns            string
		selectedName  string
		selectedKey   string
		selectedValue string

		err bool
	}{
		{
			ns:            "ns1",
			selectedName:  "secret",
			selectedKey:   "proxyA",
			selectedValue: "proxyA",
			err:           false,
		},
		{
			// Wrong selected name.
			ns:            "ns1",
			selectedName:  "proxyA",
			selectedKey:   "proxyA",
			selectedValue: "proxyA",
			err:           true,
		},
		{
			// Wrong namespace.
			ns:            "ns2",
			selectedName:  "secret",
			selectedKey:   "proxyA",
			selectedValue: "proxyA",
			err:           true,
		},
		{
			// Wrong not found selected key.
			ns:            "ns1",
			selectedName:  "secret",
			selectedKey:   "proxyD",
			selectedValue: "proxyD",
			err:           true,
		},
	} {

		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			proxyConfig := monitoringv1.ProxyConfig{
				ProxyConnectHeader: map[string][]v1.SecretKeySelector{
					"header": {
						{
							LocalObjectReference: v1.LocalObjectReference{
								Name: tc.selectedName,
							},
							Key: tc.selectedKey,
						},
					},
				},
			}

			err := store.AddProxyConfig(context.Background(), tc.ns, proxyConfig)

			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			b, err := store.ForNamespace(tc.ns).GetSecretKey(proxyConfig.ProxyConnectHeader["header"][0])
			require.NoError(t, err)
			require.Equal(t, tc.selectedValue, string(b))
		})
	}

}

func TestAddTLSConfig(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm",
				Namespace: "ns1",
			},
			Data: map[string]string{
				"cmCA":   caPEM,
				"cmCert": certPEM,
				"cmKey":  keyPEM,
			},
		},
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"secretCA":   []byte(caPEM),
				"secretCert": []byte(certPEM),
				"secretKey":  []byte(keyPEM),

				"invalidCA": []byte("invalidCA"),
				"wrongKey":  []byte("wrongKey"),
			},
		},
	)

	for _, tc := range []struct {
		ns        string
		tlsConfig *monitoringv1.TLSConfig

		err          bool
		expectedCA   string
		expectedCert string
		expectedKey  string
	}{
		{
			// CA, cert and key in secret.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			expectedCA:   caPEM,
			expectedCert: certPEM,
			expectedKey:  keyPEM,
		},
		{
			// CA in configmap, cert and key in secret.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			expectedCA:   caPEM,
			expectedCert: certPEM,
			expectedKey:  keyPEM,
		},
		{
			// CA and cert in configmap, key in secret.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			expectedCA:   caPEM,
			expectedCert: certPEM,
			expectedKey:  keyPEM,
		},
		{
			// Wrong namespace.
			ns: "ns2",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Wrong configmap selector for CA.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "secretCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Wrong secret selector for CA.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Wrong configmap selector for cert.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Wrong secret selector for cert.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						ConfigMap: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "cm",
							},
							Key: "cmCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "cmCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Wrong key selector.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "cmKey",
					},
				},
			},

			err: true,
		},
		{
			// Cert without key.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
				},
			},

			err: true,
		},
		{
			// Key without cert.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCA",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
		{
			// Cert with wrong key.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "wrongKey",
					},
				},
			},

			err: true,
		},
		{
			// Invalid CA certificate.
			ns: "ns1",
			tlsConfig: &monitoringv1.TLSConfig{
				SafeTLSConfig: monitoringv1.SafeTLSConfig{
					CA: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "invalidCA",
						},
					},
					Cert: monitoringv1.SecretOrConfigMap{
						Secret: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: "secret",
							},
							Key: "secretCert",
						},
					},
					KeySecret: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "secret",
						},
						Key: "secretKey",
					},
				},
			},

			err: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			err := store.AddSafeTLSConfig(context.Background(), tc.ns, &tc.tlsConfig.SafeTLSConfig)
			if tc.err {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			tlsAssets := store.TLSAssets()

			key := tlsAssetKeyFromSelector(tc.ns, tc.tlsConfig.CA).toString()
			b, found := tlsAssets[key]
			require.True(t, found)
			require.Equal(t, tc.expectedCA, string(b))

			key = tlsAssetKeyFromSelector(tc.ns, tc.tlsConfig.Cert).toString()
			b, found = tlsAssets[key]
			require.True(t, found)
			require.Equal(t, tc.expectedCert, string(b))

			key = tlsAssetKeyFromSecretSelector(tc.ns, tc.tlsConfig.KeySecret).toString()
			b, found = tlsAssets[key]
			require.True(t, found)
			require.Equal(t, tc.expectedKey, string(b))
		})
	}
}

func TestAddAuthorization(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
	)

	for _, tc := range []struct {
		ns           string
		selectedName string
		selectedKey  string
		authType     string

		err      bool
		expected string
	}{
		{
			ns:           "ns1",
			selectedName: "secret",
			selectedKey:  "key1",
			authType:     "Bearer",

			expected: "val1",
		},
		{
			ns:           "ns1",
			selectedName: "secret",
			selectedKey:  "key1",
			authType:     "Token",

			expected: "val1",
		},
		{
			ns:           "ns1",
			selectedName: "secreet",
			selectedKey:  "key1",
			authType:     "Token",

			err: true,
		},
		{
			ns:           "ns1",
			selectedName: "",
			selectedKey:  "",
			authType:     "Bearer",

			expected: "",
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			sel := &monitoringv1.Authorization{
				SafeAuthorization: monitoringv1.SafeAuthorization{
					Type: tc.authType,
					Credentials: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: tc.selectedName},
						Key: tc.selectedKey,
					},
				},
			}

			err := store.AddAuthorizationCredentials(context.Background(), tc.ns, sel)

			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if sel.Credentials.Name == "" {
				return
			}

			b, err := store.ForNamespace(tc.ns).GetSecretKey(*sel.Credentials)
			require.NoError(t, err)

			s := string(b)
			require.Equal(t, tc.expected, s, "expecting %q, got %q", tc.expected, s)
		})
	}
}

func TestAddAuthorizationNoCredentials(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
	)

	t.Run("", func(t *testing.T) {
		store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

		sel := &monitoringv1.Authorization{
			SafeAuthorization: monitoringv1.SafeAuthorization{
				Type: "authType",
			},
			CredentialsFile: "/path/to/secret",
		}

		err := store.AddAuthorizationCredentials(context.Background(), "foo", sel)
		require.NoError(t, err)
	})
}

func TestAddSigV4(t *testing.T) {
	const (
		accessKey = "accessKey"
		secretKey = "secretKey"
	)
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				accessKey: []byte("val1"),
				secretKey: []byte("val2"),
			},
		},
	)

	for _, tc := range []struct {
		title                string
		ns                   string
		selectedName         string
		accessKey, secretKey string

		err                 bool
		expectedAccessKeyID string
		expectedSecretKeyID string
	}{
		{
			title:        "valid access and secret keys",
			ns:           "ns1",
			selectedName: "secret",
			accessKey:    accessKey,
			secretKey:    secretKey,

			expectedAccessKeyID: "val1",
			expectedSecretKeyID: "val2",
		},
		{
			title:        "wrong namespace",
			ns:           "ns2",
			selectedName: "secret",
			accessKey:    accessKey,
			secretKey:    secretKey,

			err: true,
		},
		{
			title:        "wrong name",
			ns:           "ns1",
			selectedName: "faulty",
			accessKey:    accessKey,
			secretKey:    secretKey,

			err: true,
		},
		{
			title:        "wrong key selector",
			ns:           "ns1",
			selectedName: "secret",
			accessKey:    "wrong-access-key",
			secretKey:    "wrong-secret-key",

			err: true,
		},
		{
			title:        "missing access key",
			ns:           "ns1",
			selectedName: "secret",
			secretKey:    secretKey,

			err: true,
		},
		{
			title:        "missing secret key",
			ns:           "ns1",
			selectedName: "secret",
			accessKey:    accessKey,

			err: true,
		},
		{
			title:        "empty keys",
			ns:           "ns1",
			selectedName: "secret",
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			sigV4 := monitoringv1.Sigv4{}
			if tc.accessKey != "" {
				sigV4.AccessKey = &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: tc.selectedName,
					},
					Key: tc.accessKey,
				}
			}
			if tc.secretKey != "" {
				sigV4.SecretKey = &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: tc.selectedName,
					},
					Key: tc.secretKey,
				}
			}

			err := store.AddSigV4(context.Background(), tc.ns, &sigV4)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			if sigV4.AccessKey != nil {
				b, err := store.ForNamespace(tc.ns).GetSecretKey(*sigV4.AccessKey)
				require.NoError(t, err)
				require.Equal(t, tc.expectedAccessKeyID, string(b))
			}

			if sigV4.SecretKey != nil {
				b, err := store.ForNamespace(tc.ns).GetSecretKey(*sigV4.SecretKey)
				require.NoError(t, err)
				require.Equal(t, tc.expectedSecretKeyID, string(b))
			}
		})
	}
}

func TestAddAzureOAuth(t *testing.T) {
	const (
		clientSecret = "clientSecretKey"
	)
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				clientSecret: []byte("val1"),
			},
		},
	)

	for _, tc := range []struct {
		title                string
		ns                   string
		selectedName         string
		accessKey, secretKey string

		err      bool
		expected string
	}{
		{
			title:        "valid clientSecret key",
			ns:           "ns1",
			selectedName: "secret",
			secretKey:    clientSecret,

			expected: "val1",
		},
		{
			title:        "wrong namespace",
			ns:           "ns2",
			selectedName: "secret",
			secretKey:    clientSecret,

			err: true,
		},
		{
			title:        "wrong name",
			ns:           "ns1",
			selectedName: "faulty",
			secretKey:    clientSecret,

			err: true,
		},
		{
			title:        "wrong key selector",
			ns:           "ns1",
			selectedName: "secret",
			secretKey:    "wrong-secret-key",

			err: true,
		},
	} {
		t.Run("", func(t *testing.T) {
			store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

			azureAD := monitoringv1.AzureAD{}
			azureOAuth := monitoringv1.AzureOAuth{}
			if tc.secretKey != "" {
				azureOAuth.ClientSecret = v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: tc.selectedName,
					},
					Key: tc.secretKey,
				}
			}
			azureAD.OAuth = &azureOAuth

			err := store.AddAzureOAuth(context.Background(), tc.ns, &azureAD)
			if tc.err {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			b, err := store.ForNamespace(tc.ns).GetSecretKey(azureOAuth.ClientSecret)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(b))
		})
	}
}

func TestUpdateObject(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
	)
	store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

	// Add the secret to the store by fetching it
	sel := v1.SecretKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "secret",
		},
		Key: "key1",
	}
	val, err := store.GetSecretKey(context.Background(), "ns1", sel)
	require.NoError(t, err)
	require.Equal(t, "val1", val)

	// Update the secret object
	updatedSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "ns1",
		},
		Data: map[string][]byte{
			"key1": []byte("val2"),
		},
	}
	err = store.UpdateObject(updatedSecret)
	require.NoError(t, err)

	// Now, getting the key should return the updated value
	val, err = store.GetSecretKey(context.Background(), "ns1", sel)
	require.NoError(t, err)
	require.Equal(t, "val2", val)

	// Test updating with nil object
	err = store.UpdateObject(nil)
	require.Error(t, err)
}

func TestDeleteObject(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm",
				Namespace: "ns1",
			},
			Data: map[string]string{
				"cmKey": "cmVal",
			},
		},
	)
	store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

	// Add secret and configmap to the store by fetching them
	secretSel := v1.SecretKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "secret",
		},
		Key: "key1",
	}
	val, err := store.GetSecretKey(context.Background(), "ns1", secretSel)
	require.NoError(t, err)
	require.Equal(t, "val1", val)

	cmSel := v1.ConfigMapKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "cm",
		},
		Key: "cmKey",
	}
	_, err = store.GetConfigMapKey(context.Background(), "ns1", cmSel)
	require.NoError(t, err)

	// Try deleting the secret object
	secretObj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "ns1",
		},
	}
	err = store.DeleteObject(secretObj)
	require.NoError(t, err)

	// Also delete the secret from the fake clientset to simulate full removal
	err = c.CoreV1().Secrets("ns1").Delete(context.Background(), "secret", metav1.DeleteOptions{})
	require.NoError(t, err)

	// Now, getting the key should fail since the secret is deleted from the store
	_, err = store.GetSecretKey(context.Background(), "ns1", secretSel)
	require.Error(t, err)

	// Try deleting the configmap object (should not error even if it doesn't exist in the client)
	cmObj := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: "ns1",
		},
	}
	err = store.DeleteObject(cmObj)
	require.NoError(t, err)

	err = c.CoreV1().ConfigMaps("ns1").Delete(context.Background(), "cm", metav1.DeleteOptions{})
	require.NoError(t, err)

	// Now, getting the key should fail since the configmap is deleted from the store
	_, err = store.GetConfigMapKey(context.Background(), "ns1", cmSel)
	require.Error(t, err)

	// Test deleting with nil object
	err = store.DeleteObject(nil)
	require.Error(t, err)
}

func TestGetObject(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm",
				Namespace: "ns1",
			},
			Data: map[string]string{
				"cmKey": "cmVal",
			},
		},
	)
	store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

	// Add secret and configmap to the store by fetching them
	secretSel := v1.SecretKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "secret",
		},
		Key: "key1",
	}
	_, err := store.GetSecretKey(context.Background(), "ns1", secretSel)
	require.NoError(t, err)

	cmSel := v1.ConfigMapKeySelector{
		LocalObjectReference: v1.LocalObjectReference{
			Name: "cm",
		},
		Key: "cmKey",
	}
	_, err = store.GetConfigMapKey(context.Background(), "ns1", cmSel)
	require.NoError(t, err)

	// Test getting existing secret
	secretObj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "ns1",
		},
	}
	obj, exists, err := store.GetObject(secretObj)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, obj)
	secret, ok := obj.(*v1.Secret)
	require.True(t, ok)
	require.Equal(t, "secret", secret.Name)
	require.Equal(t, "ns1", secret.Namespace)
	require.Equal(t, []byte("val1"), secret.Data["key1"])

	// Test getting existing configmap
	cmObj := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm",
			Namespace: "ns1",
		},
	}
	obj, exists, err = store.GetObject(cmObj)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, obj)
	cm, ok := obj.(*v1.ConfigMap)
	require.True(t, ok)
	require.Equal(t, "cm", cm.Name)
	require.Equal(t, "ns1", cm.Namespace)
	require.Equal(t, "cmVal", cm.Data["cmKey"])

	// Test getting non-existing object
	nonExistingSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "notfound",
			Namespace: "ns1",
		},
	}
	obj, exists, err = store.GetObject(nonExistingSecret)
	require.NoError(t, err)
	require.False(t, exists)
	require.Nil(t, obj)

	// Test getting with nil object
	obj, exists, err = store.GetObject(nil)
	require.Error(t, err)
	require.False(t, exists)
	require.Nil(t, obj)
}

func TestAddObject(t *testing.T) {
	c := fake.NewSimpleClientset(
		&v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret",
				Namespace: "ns1",
			},
			Data: map[string][]byte{
				"key1": []byte("val1"),
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "cm",
				Namespace: "ns1",
			},
			Data: map[string]string{
				"cmKey": "cmVal",
			},
		},
	)
	store := NewStoreBuilder(c.CoreV1(), c.CoreV1())

	// Add a secret object
	secretObj := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret2",
			Namespace: "ns1",
		},
		Data: map[string][]byte{
			"key2": []byte("val2"),
		},
	}
	err := store.AddObject(secretObj)
	require.NoError(t, err)

	// Retrieve the secret object
	obj, exists, err := store.GetObject(secretObj)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, obj)
	secret, ok := obj.(*v1.Secret)
	require.True(t, ok)
	require.Equal(t, "secret2", secret.Name)
	require.Equal(t, "ns1", secret.Namespace)
	require.Equal(t, []byte("val2"), secret.Data["key2"])

	// Add a configmap object
	cmObj := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cm2",
			Namespace: "ns1",
		},
		Data: map[string]string{
			"cmKey2": "cmVal2",
		},
	}
	err = store.AddObject(cmObj)
	require.NoError(t, err)

	// Retrieve the configmap object
	obj, exists, err = store.GetObject(cmObj)
	require.NoError(t, err)
	require.True(t, exists)
	require.NotNil(t, obj)
	cm, ok := obj.(*v1.ConfigMap)
	require.True(t, ok)
	require.Equal(t, "cm2", cm.Name)
	require.Equal(t, "ns1", cm.Namespace)
	require.Equal(t, "cmVal2", cm.Data["cmKey2"])

	// Add nil object should error
	err = store.AddObject(nil)
	require.Error(t, err)
}
