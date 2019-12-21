# Stripe payment processing flow

1. Create Order,
```text
POST {{dev}}/orders
```

2. Create ClientToken,
```text
POST {{dev}}/payment-gateways/token
```

3. Initiate Payment with callback from client with,
```text
callback_url: {{dev}}/orders/{{order_id}}/pay
```

4. If payment successful, BrainTree will send the callback to server.
Once server gets the callback from BrainTree with Nonce, Server will call BrainTree for payment settlement.
Thus the payment flow completes.
