# A simple dynamic DNS tool for Cloudflare in Go <br /> 

Reference from [Cloudflare API](https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-update-dns-record) <br />

Environment variables needed:
| Environment variables | Description |
|---|---|
| `CF_ZONE_ID` | The zone ID of your domain provided by Cloudflare |
| `CF_DNS_API_KEY` | A API token provided by Cloudflare which can be obtained under the profile setting from Cloudflare dashboard |
| `CF_DNS_RECORD_ID` | The unique ID of your DNS record provided by Cloudflare <br /> Get by `curl -X GET "https://api.cloudflare.com/client/v4/zones/{zoneId}/dns\_records" -H "Authorization: Bearer {token}" -H "Content-Type:application/json"` |
| `CF_DNS_API_EMAIL` | The email from Cloudflare account where your domain belongs to |
| `TARGET_DOMAIN` | The domain name which it's record needs to be update |

![Go Gopher icon](https://github.com/doz-8108/dynamic-dns-home/assets/40817247/d127188e-d451-40fc-b414-1e08fc8755cf)
