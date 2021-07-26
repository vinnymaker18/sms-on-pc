package xyz.vinayemani.texttopc;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.telephony.SmsMessage;
import android.util.Log;

import java.io.IOException;

import static android.provider.Telephony.Sms.Intents.SMS_RECEIVED_ACTION;
import static android.provider.Telephony.Sms.Intents.getMessagesFromIntent;

public class NewSMSReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        Log.d(NEW_SMS_ON_RECEIVE_TAG, intent.getAction());

        if (SMS_RECEIVED_ACTION.equals(intent.getAction())) {
            SmsMessage[] messages = getMessagesFromIntent(intent);
            for (final SmsMessage message : messages) {
                String displayMsgBody = message.getDisplayMessageBody();
                String originAddress = message.getDisplayOriginatingAddress();
                String msgBody = message.getMessageBody();

                Thread thread = new Thread(new Runnable() {
                    @Override
                    public void run() {
                        try {
                            smsClient.postNewSms(message);
                        } catch (IOException e) {}
                    }
                });

                thread.start();

                Log.d(NEW_SMS_ON_RECEIVE_TAG, displayMsgBody);
                Log.d(NEW_SMS_ON_RECEIVE_TAG, originAddress);
                Log.d(NEW_SMS_ON_RECEIVE_TAG, msgBody);

                try {
                    thread.join();
                } catch (InterruptedException ex) {}
            }
        }
    }

    private static final String SERVER_URL = "https://minibox.vinayemani.xyz/sms";
    private static SmsClient smsClient = new SmsClient(SERVER_URL);
    private static final String NEW_SMS_ON_RECEIVE_TAG = "newSmsReceiverOnReceive";
}
