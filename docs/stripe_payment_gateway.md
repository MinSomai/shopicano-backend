# Stripe payment processing flow

1. Create Order,
```text
POST {{dev}}/orders
```

2. Generate payment nonce,
```text
POST {{dev}}/orders/{{order_id}}/nonce
```

3. Initiate payment with nonce

4. If payment successful, Stripe will send the callback to server.
Once server gets the callback from Stripe as success, Server will mark the order as `payment_completed` else `payment_failed`.
Thus the payment flow completes.
