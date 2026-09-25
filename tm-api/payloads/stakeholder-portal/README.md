# Stakeholder Portal Payloads

This directory is for generated verification payloads.

The Swagger-ordered payload reference is:

```text
payloads/stakeholder-portal/SWAGGER_ORDER_PAYLOADS.md
```

Use that file when testing from Swagger. It lists every endpoint in Swagger order, with the exact JSON body or `No JSON body`.

Run Step 0 and Step 2 in `STAKEHOLDER_PORTAL_TEST_VERIFICATION_FLOW.md` from the repo root:

```bash
cd /Users/toluwase/GolandProjects/admin-dashboard-api
mkdir -p evidence/stakeholder-portal payloads/stakeholder-portal
```

The JSON files are generated from real AdminDB and Wallet DB values. Do not hand-type random IDs into these payloads.
