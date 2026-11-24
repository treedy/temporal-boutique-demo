# Checkout `PlaceOrder` call flow

The swim lane diagram below follows the calls made by `PlaceOrder` in
[`src/checkoutservice/main.go`](/home/tawdball/workspaces/temporal-boutique-demo/src/checkoutservice/main.go).
Calls between services are labeled **gRPC**.

``` mermaid
sequenceDiagram
    autonumber
    participant FrontEnd as Front End
    participant Checkout as CheckoutService
    participant Cart as CartService
    participant Catalog as ProductCatalogService
    participant Currency as CurrencyService
    participant Shipping as ShippingService
    participant Payment as PaymentService
    participant Email as EmailService

    FrontEnd->>Checkout: PlaceOrder(request)
    Checkout->>Checkout: uuid.NewUUID()

    rect rgb(100, 100, 100)
      note right of FrontEnd: prepareOrderItemsAndShippingQuoteFromCart
      %% Checkout->>Checkout: getUserCart()
      Checkout->>Cart: GetCart(userID) [gRPC]
      Cart-->>Checkout: []cart items

        %% Checkout->>Checkout: prepOrderItems()
        loop For each cart item (prepOrderItems)
            Checkout->>Catalog: GetProduct(productID) [gRPC]
            Catalog-->>Checkout: product
            %% Checkout->>Checkout: convertCurrency(product price)
            Checkout->>Currency: Convert(from USD, to user currency) [gRPC]
            Currency-->>Checkout: localized price
        end

        %% Checkout->>Checkout: quoteShipping()
        Checkout->>Shipping: GetQuote(address, cart items) [gRPC]
        Shipping-->>Checkout: USD shipping cost
        %% Checkout->>Checkout: convertCurrency(shipping cost)
        Checkout->>Currency: Convert(from USD, to user currency) [gRPC]
        Currency-->>Checkout: localized shipping cost
        Checkout-->>Checkout: order preparation result
    end

    Checkout->>Checkout: calculate order total
    %% Checkout->>Checkout: chargeCard()
    Checkout->>Payment: Charge(total, credit card) [gRPC]
    Payment-->>Checkout: transaction ID

    %% Checkout->>Checkout: shipOrder()
    Checkout->>Shipping: ShipOrder(address, cart items) [gRPC]
    Shipping-->>Checkout: tracking ID

    %% Checkout->>Checkout: emptyUserCart()
    Checkout->>Cart: EmptyCart(userID) [gRPC]
    Cart-->>Checkout: empty-cart result

    Checkout->>Checkout: build OrderResult
    %%Checkout->>Checkout: sendOrderConfirmation()
    Checkout->>Email: SendOrderConfirmation(email, order) [gRPC]
    Email-->>Checkout: confirmation result
    Checkout-->>FrontEnd: PlaceOrderResponse
```

Original Diagram
```mermaid
sequenceDiagram
    autonumber
    participant Caller
    participant Checkout as CheckoutService
    participant Cart as CartService
    participant Catalog as ProductCatalogService
    participant Currency as CurrencyService
    participant Shipping as ShippingService
    participant Payment as PaymentService
    participant Email as EmailService

    Caller->>Checkout: PlaceOrder(request)
    Checkout->>Checkout: uuid.NewUUID()
    Checkout->>Checkout: prepareOrderItemsAndShippingQuoteFromCart()

    Checkout->>Checkout: getUserCart()
    Checkout->>Cart: GetCart(userID) [gRPC]
    Cart-->>Checkout: cart items

    Checkout->>Checkout: prepOrderItems()
    loop For each cart item
        Checkout->>Catalog: GetProduct(productID) [gRPC]
        Catalog-->>Checkout: product
        Checkout->>Checkout: convertCurrency(product price)
        Checkout->>Currency: Convert(from USD, to user currency) [gRPC]
        Currency-->>Checkout: localized price
    end

    Checkout->>Checkout: quoteShipping()
    Checkout->>Shipping: GetQuote(address, cart items) [gRPC]
    Shipping-->>Checkout: USD shipping cost
    Checkout->>Checkout: convertCurrency(shipping cost)
    Checkout->>Currency: Convert(from USD, to user currency) [gRPC]
    Currency-->>Checkout: localized shipping cost
    Checkout-->>Checkout: order preparation result

    Checkout->>Checkout: calculate order total
    Checkout->>Checkout: chargeCard()
    Checkout->>Payment: Charge(total, credit card) [gRPC]
    Payment-->>Checkout: transaction ID

    Checkout->>Checkout: shipOrder()
    Checkout->>Shipping: ShipOrder(address, cart items) [gRPC]
    Shipping-->>Checkout: tracking ID

    Checkout->>Checkout: emptyUserCart()
    Checkout->>Cart: EmptyCart(userID) [gRPC]
    Cart-->>Checkout: empty-cart result

    Checkout->>Checkout: build OrderResult
    Checkout->>Checkout: sendOrderConfirmation()
    Checkout->>Email: SendOrderConfirmation(email, order) [gRPC]
    Email-->>Checkout: confirmation result
    Checkout-->>Caller: PlaceOrderResponse
```
