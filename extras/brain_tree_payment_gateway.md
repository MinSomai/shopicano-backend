# BrainTree payment processing flow

1. Create Order,
```text
POST {{dev}}/orders
```

2. Get ClientToken,
```text
GET {{dev}}/payments/configs
```

3. Initiate Payment with callbackFunction using brainTree client and on callback function you will get payload with nonce. Send nonce to server.
```text
POST {{dev}}/orders/{{order_id}}/pay

{
    "nonce": "tokencc_bj_71nryc_6qbzqs_rqbks6_fj2s87_at3"
}
```

4. If payment is successful, you will get success response from Server.
Thus the payment flow completes.
