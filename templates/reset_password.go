package templates

import (
	"bytes"
	"html/template"
)

var resetPasswordTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1">

    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/inter-ui@3.12.0/inter.min.css">

    <title>{{ .platformName }} | Reset Password</title>

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
        padding: 48px 47px;
    }
    .my-28 {
        margin-top: 28px;
        margin-bottom: 28px;
    }
    .btn{
        width: 100%;
        border-radius: 3px;
        background-color: #3f71f4;
        padding-top: 21px;
        padding-bottom: 21px;
        font-size: 14px;
        font-weight: bold;
        font-stretch: normal;
        font-style: normal;
        line-height: normal;
        letter-spacing: 0.53px;
        text-align: center;
        color: #ffffff;
        font-weight: 400;
        vertical-align: middle;
        cursor: pointer;
        -webkit-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
        user-select: none;
        border: 1px solid transparent;
    }
    cp{
        font-size: 14px;
        font-weight: normal;
        font-stretch: normal;
        font-style: normal;
        line-height: 1.29;
        letter-spacing: normal;
        color: #6b7694;
        text-align: center!important;
    }
    </style>
</head>


<body>
    <center>
        <table width="100%" border="0" cellspacing="0" cellpadding="0" style="margin: 0; padding-top: 138px; width: 100%; height: 100%;">
            <tr>
                <td style="margin: 0; padding: 0; width: 100%; height: 100%;" align="center">
                    <a href="{{ .platformWebsite }}" target="_blank"><img src="group-26@3x.png" width="165px" height="42px" alt=""></a>
                    <table width="600" border="0" cellspacing="0" cellpadding="0" style="margin-top: 38px; padding: 0;">
                        <tr>
                            <td class="container" style="width:600px; min-width:600px; width: 100%;" align="center">
                                <img src="key.png" width="74" height="74" alt="">
                                <h3 class="my-28">Reset Your Password</h3>

                                <p>
                                    Please click the button below to reset your password. This is valid for 24 hours, 
                                    then you have to request new one.
                                </p>

                                <button class="btn">Reset Password</button>

                                <p class="my-28">
                                    If you’re having trouble with the button ‘Reset Password', 
                                    copy and paste the URL below into your web browser.
                                </p>

                                <a href="{{ .resetPasswordUrl }}">{{ .resetPasswordUrl }}</a>
                            </td>
                        </tr>
                    </table>

                    <table width="600" border="0" cellspacing="0" cellpadding="0" style="margin-top: 0; padding: 0;">
                        <tr>
                            <td style="width:600px; min-width:600px; width: 100%;" align="center">
                                <p style="font-size: 14px;font-weight: normal;font-stretch: normal;font-style: normal;line-height: 1.29;letter-spacing: normal;color: #6b7694;text-align: center!important">
                                    © 2020 {{ .platformName }}. All rights reserved.
                                </P>
                                <p style="font-size: 14px;font-weight: normal;font-stretch: normal;font-style: normal;line-height: 1.29;letter-spacing: normal;color: #6b7694;text-align: center!important">
                                    Powered by <a href="{{ .platformWebsite }}" target="_blank" style="text-decoration: none">{{ .platformName }}</a>
                                </P>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </center>
</body>
</html>
`

func GenerateResetPasswordEmailHTML(params map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	t := template.Must(template.New("ResetPasswordTemplate").Parse(resetPasswordTemplate))
	if err := t.Execute(&buf, params); err != nil {
		return "", err
	}
	return buf.String(), nil
}
