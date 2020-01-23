var aws = require('aws-sdk');
var ses = new aws.SES();

/*
    Developed by "https://github.com/vitorfmc"
    
    =======================================================
    Overview:
    =======================================================

    This Lambda Function is example of integration with SES Service and DynamoDB.
    The idea is: Everytime a table in DynamoDB receive a data, it will send the event 
    information to this lambda function, which will send a e-mail for registration.

    DynamoDB Stream ==> This Lambda Function ==> SES

    Obs.: Remember to give Dynamo and SES policies to your Lambda Function
*/

exports.handler = (event, context, callback) => {
    
     /* Obs.: The "Source" is configured at AWS SES panel */
     var params = {
        Destination: {
            ToAddresses: ["????@gmail.com"]
        },
        Message: {
            Body: {
                Text: { 
                  Data: "Test"
                }
                
            },
            Subject: { 
              Data: "Test Email"
            }
        },
        Source:"????@gmail.com"
    };

    
     ses.sendEmail(params, function (err, data) {
        callback(null, {err: err, data: data});
        if (err) {
            console.log(err);
            context.fail(err);
        } else {
            console.log(data);
            context.succeed(event);
        }
    });
};