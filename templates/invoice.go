package templates

import (
	"bytes"
	"html/template"
)

var invoiceTemplate = `
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1">

    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/inter-ui@3.12.0/inter.min.css">

    <title>{{ .platformName}} | Order Invoice</title>

    <style type="text/css" media="screen">
    body { padding:0 !important; margin:0 auto !important; font-family: Inter; display:block !important; min-width:100% !important; width:100% !important; background: #f6f8fc;; -webkit-text-size-adjust:none }

    p {
        font-size: 16px;
        font-weight: normal;
        font-stretch: normal;
        font-style: normal;
        line-height: 1.5;
        letter-spacing: normal;
        color: #5a637c;
        text-align: center;
    }
    a{color: #3f71f4; word-break: break-all; text-align: left;}
    h3{
        font-size: 24px;
        font-weight: 500;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: normal;
        text-align: center;
        color: #363b4a;
    }
    img { position: relative; margin: 0 !important; -ms-interpolation-mode: bicubic;}


    .container{
        border-radius: 3px;
        box-shadow: -2px -3px 8px 0 rgba(255, 255, 255, 0.5);
        border: solid 1px #e9eceb;
        background-color: #ffffff;
        padding: 44px 28px;
    }
    .title {
        margin-top: 40px;
        margin-bottom: 16px;
        padding: 4px 1px;
    }
    .wrapper{
        display: block;
        height: auto;
        width: 100%;
        position: relative;
        text-align: left;
        margin-bottom: 46px;
        margin-top: 32px;
    }
    .details{
        font-size: 12px;
        font-weight: 600;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: normal;
        color: #363b4a;
        text-align: left;
    }
    .address{
        font-size: 13px;
        font-weight: normal;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: normal;
        color: #3f4a5a;
        text-align: left;
    }
    .tbl-heder{
        font-size: 12px;
        font-weight: 600;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: normal;
        text-align: right;
        color: #363b4a;
    }
    .tbl-heder th{
        padding: 8px 0;
    }
    .tbl-data{
        font-size: 13px;
        font-weight: normal;
        font-stretch: normal;
        font-style: normal;
        line-height: 1.38;
        letter-spacing: normal;
        text-align: right;
        color: #363b4a;
        text-align: left;
    }
    .td-border1st{
        border: 1px solid #ecf1fe; 
        border-left: 0; 
        border-right: 0; 
        padding-top: 9px; 
        padding-bottom: 12px;
    }
    .td-border2nd, .border_top{
        border: 1px solid #ecf1fe; 
        border-left: 0; 
        border-right: 0; 
        border-bottom: 0;
        padding: 7px 0;
    }
    .border_bottom{
        border: 1px solid #ecf1fe; 
        border-left: 0; 
        border-right: 0; 
        border-top: 0;
    }
    .footer{
        position: relative;
    }
    .footer ul{
        margin: 0;
        padding: 0;
        display: inline-block;
    }
    .footer ul li{
        list-style: none;
        display: inline-block;
        padding: 0;
        margin: 0;
        padding-right: 14px;
        font-size: 13px;
        font-weight: normal;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: normal;
        color: #8993a4;
    }
    .footer ul li a{
        color: #8993a4;
        text-decoration: none;
    }
    .oval{
        width: 23px;
        height: 23px;
        border-radius: 15px;
        background-color: #c2d2ff;
    }
    </style>
</head>


<body>
    <center>
        <table width="100%" border="0" cellspacing="0" cellpadding="0" style="margin: 0; padding-top: 138px; padding-bottom: 138px; width: 100%; height: 100%;">
            <tr>
                <td style="margin: 0; padding: 0; width: 100%; height: 100%;" align="center">
                    <table width="600" border="0" cellspacing="0" cellpadding="0" style="margin-top: 38px; padding: 0;">
                        <tr>
                            <td class="container" style="width:600px; padding-bottom: 92px; min-width:600px; width: 100%;" align="center">
                                <img src="group-26@3x.png" width="165px" height="42px" alt="">
                                <h3 class="title">{{ .greetings }}</h3>

                                <p>{{ .intros }}</p>

                                <div class="wrapper">
                                    <p class="details">Billing Details</p>
                                    <p class="address">{{ .billingAddress }}</p>
                                </div>

                                <!--    1st table   -->
                                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                                    <tr class="tbl-heder" style="text-align: left;">
                                        <th>Order ID</th>
                                        <th>BUYER</th> 
                                        <th style="text-align: right;">INVOICE DATE</th>
                                    </tr>
                                    <tr class="tbl-data">
                                        <td class="td-border1st" style="color: #3f71f4;"><a href="{{ .orderUrl }}" target="_blank">#{{ .orderHash }}</a></td>
                                        <td class="td-border1st">{{ .buyerName }}</td>
                                        <td class="td-border1st" style="text-align: right;">{{ .orderDate }}</td>
                                    </tr>
                                </table>
                                <!--    1st table end   -->

                                <!--    2nd table   -->
                                <table width="100%" border="0" cellspacing="0" cellpadding="0">
                                    <tr class="tbl-heder" style="text-align: left;">
                                        <th>DESCRIPTION</th>
                                        <th style="text-align: center;">Price</th> 
                                        <th style="text-align: center;">Qty</th> 
                                        <th style="text-align: right;">Sub Total</th>
                                    </tr>
                                    {{ range $item := .orderedItems }}
                                        <tr class="tbl-data">
                                            <td class="td-border2nd" style="padding: 7px 0;">{{ $item.name }}</td>
                                            <td class="td-border2nd" style="text-align: center; padding: 7px 0;">{{ $item.price }}</td>
                                            <td class="td-border2nd" style="text-align: center; padding: 7px 0;">{{ $item.quantity }}</td>
                                            <td class="td-border2nd" style="text-align: right; padding: 7px 0;">{{ $item.subTotal }}</td>
                                        </tr>
                                    {{end}}

                                    <tr class="tbl-data">
                                        <td class="border_top" style="padding: 7px 0;"></td>
                                        <td class="border_top" style="text-align: center;"></td>
                                        <td class="border_top" style="text-align: right; padding: 10px 0 6px 0;">Subtotal:</td>
                                        <td class="border_top" style="text-align: right; padding: 10px 0 6px 0;">{{ .subTotal }}</td>
                                    </tr>

                                    {{ if .isCouponApplied }}
                                        <tr class="tbl-data">
                                            <td class="border_bottom" style="padding: 7px 0;"></td>
                                            <td class="border_bottom" style="text-align: center;"></td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">Discount:</td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">{{ .discount }}</td>
                                        </tr>
                                    {{end}}

                                    {{ if .isCouponApplied }}
                                        <tr class="tbl-data">
                                            <td class="border_bottom" style="padding: 7px 0;"></td>
                                            <td class="border_bottom" style="text-align: center;"></td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">Coupon Code:</td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">{{ .couponCode }}</td>
                                        </tr>
                                    {{end}}

                                    {{ if .isShippable }}
                                        <tr class="tbl-data">
                                            <td class="border_bottom" style="padding: 7px 0;"></td>
                                            <td class="border_bottom" style="text-align: center;"></td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">Shipping Charge:</td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">{{ .shippingCharge }}</td>
                                        </tr>
                                    {{end}}

                                    {{ if .isDigitalPayment }}
                                        <tr class="tbl-data">
                                            <td class="border_bottom" style="padding: 7px 0;"></td>
                                            <td class="border_bottom" style="text-align: center;"></td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">Payment Processing Fee:</td>
                                            <td class="border_bottom" style="text-align: right; padding: 2px 0 10px 0;">{{ .paymentProcessingFee }}</td>
                                        </tr>
                                    {{end}}

                                    <tr class="tbl-data">
                                        <th style="text-align: left;">Shipping Address</th>
                                        <td style="text-align: center;"></td>
                                        <td style="text-align: right; padding: 13px 0;">Total:</td>
                                        <td style="text-align: right; padding: 13px 0;">{{ .grandTotal }}</td>
                                    </tr>
                                    <tr class="tbl-data">
                                        <td style="text-align: left; padding: 1px 10px 1px 0;">{{ .shippingAddress }}</td>
                                        <td style="text-align: center;"></td>
                                        <td style="text-align: right; padding: 1px 0;">Payment Status:</td>
                                        <td style="text-align: right; padding: 1px 0; color: #ff8124; font-weight: 500; font-size: 14px;">{{ .paymentStatus }}</td>
                                    </tr>
									{{ if .isShippable }}
                                    	<tr class="tbl-data">
                                        	<td style="text-align: left; padding: 1px 10px 1px 0;"></td>
                                        	<td style="text-align: center;"></td>
                                        	<td style="text-align: right; padding: 1px 0;">Order Status:</td>
                                        	<td style="text-align: right; padding: 1px 0; color: #ff8124; font-weight: 500; font-size: 14px;">{{ .status }}</td>
                                    	</tr>
									{{end}}
                                </table>
                                <!--    2nd table end  -->
                            </td>
                        </tr>

                        <tr style="width: 600px; height: 82px; border: solid 1px #e9eceb; background-color: #ecf1fe;">
                            <td class="footer" style="padding: 33px 26px;">
                                <ul>
                                    <li><a href="{{ .siteUrl }}">License</a></li>
                                    <li><a href="{{ .siteUrl }}">Policy</a></li>
                                    <li><a href="{{ .siteUrl }}">Contact us</a></li>
                                    <li><a href="{{ .siteUrl }}">Help</a></li>
                                </ul>
                                <ul style="float: right;">
                                    <li style="background-color: #c2d2ff; padding: 2px; border-radius: 20px; margin-left: 14px;">
                                        <a style="line-height: 0; padding: 0; margin: 0;" href="#">
                                            <img style="display: inline-block; margin-left: auto; margin-right: auto; width:23px; height: 23px; padding: 2px 2px 0 2px;" src="{{ .assetsUrl }}/icon_fb.png" alt="">
                                        </a>
                                    </li>
                                    <li style="background-color: #c2d2ff; padding: 2px; border-radius: 20px; margin-left: 14px;">
                                        <a style="line-height: 0; padding: 0; margin: 0;" href="#">
                                            <img style="display: inline-block; margin-left: auto; margin-right: auto; width:23px; height: 23px; padding: 2px 2px 0 2px;" src="{{ .assetsUrl }}/icon_insta.png" alt="">
                                        </a>
                                    </li>
                                    <!-- <li style="background-color: #c2d2ff; padding: 2px; border-radius: 20px; margin-left: 14px;">
                                        <a style="line-height: 0; padding: 0; margin: 0;" href="#">
                                            <img style="display: inline-block; margin-left: auto; margin-right: auto; width:23px; height: 23px; padding: 2px 2px 0 2px;" src="{{ .assetsUrl }}/icon_dribbl.png" alt="">
                                        </a>
                                    </li> -->
                                </ul>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </center>
</body>
`

func GenerateInvoiceEmailHTML(params map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("InvoiceTemplate").Parse(invoiceTemplate))
	if err := t.Execute(&buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}
