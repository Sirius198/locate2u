# Locate2u Integration

### Notes:
- Use `l.Error("test: %s", [arg1, arg2...])` (also has `Warning()`, `Info()` and `Debug()` functions) for logging - this will be replaced by a production logger later

---

## TODO:
- [x] Setup project repo
- [ ] Locate2u Authentication
- [ ] API authentication
- [ ] Sync
    - [ ] Sync customers from API to Locate2u
        - [ ] Save the customer's API id as a customer field inside Locate2u?
        - [ ] Save the customer's default shipping address as the main address in Locate2u
    - [ ] Update customers from API to Locate2u
    - [ ] Delete customers no longer from API in Locate2u
- [ ] Create http endpoint `/trip/{fulfillmentid}`
    - Lookup fulfillment from API (e.g. `IF1234`)
    - Create a `stop` in Locate2u based on that fulfillment
        - Include `name`, `lines`, `notes`, `address` from fulfillment
        - Set `tripDate` to current date
        - Set the `assignedTeamMemberId` to a value stores in config
        - Set the `customerId` by finding the correct customer synced into Locate2U
        - Set the `runNumber` to `1` if it is currently the morning and `2` if it is currently the afternoon
- [ ] Create a new 'Link' in Locate2u when new Stops are created (when the above endpoint is used)
    - Save the link to the API using the `add-tracking` endpoint (detailed below)
    - For the message in the link use a message stored in the config or default to an empty message

---

## Locate2u Api

- Working with Locate2u API: https://support.locate2u.com/article/e487ymrz1m-locate-2-u-s-rest-api
- Experimenting with Locate2u API Call: https://support.locate2u.com/article/68e5ttocqt/preview
- Swagger API reference: https://api-test.locate2u.com/index.html

---

## Our Api 

### Authentication

- Base URL: `https://api.colorex.co.nz`
- Authentication Header: `Authorization: Bearer {provided token}`

### Endpoints

- Customers
    - GET `/api/customer?page=1&per_page=100`
- Transaction by transaction number
   - GET `/api/transaction/tranid/{transactionid}`
- Add tracking link
    - GET `/ext/add-tracking/{transactionid}?delivery=true&complete=true&tracking={trackinglink}`
