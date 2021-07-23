package xyz.vinayemani.texttopc;

import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.telephony.SmsMessage;
import android.util.Log;

import static android.provider.Telephony.Sms.Intents.SMS_RECEIVED_ACTION;
import static android.provider.Telephony.Sms.Intents.getMessagesFromIntent;

public class NewSMSReceiver extends BroadcastReceiver {
    @Override
    public void onReceive(Context context, Intent intent) {
        Log.d(NEW_SMS_ON_RECEIVE_TAG, intent.getAction());
        if (SMS_RECEIVED_ACTION.equals(intent.getAction())) {
            SmsMessage[] messages = getMessagesFromIntent(intent);
            for (SmsMessage message : messages) {
                String displayMsgBody = message.getDisplayMessageBody();
                String originAddress = message.getDisplayOriginatingAddress();
                String msgBody = message.getMessageBody();
                Log.d(NEW_SMS_ON_RECEIVE_TAG, displayMsgBody);
                Log.d(NEW_SMS_ON_RECEIVE_TAG, originAddress);
                Log.d(NEW_SMS_ON_RECEIVE_TAG, msgBody);
            }
        }
    }

    private static final String NEW_SMS_ON_RECEIVE_TAG = "newSmsReceiverOnReceive";
}
