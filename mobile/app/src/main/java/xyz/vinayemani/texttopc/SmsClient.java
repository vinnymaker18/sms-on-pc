package xyz.vinayemani.texttopc;

import android.telephony.SmsMessage;
import android.util.Log;

import com.android.volley.Request;
import com.android.volley.RequestQueue;
import com.android.volley.Response;
import com.android.volley.toolbox.StringRequest;
import com.android.volley.toolbox.Volley;

import java.io.BufferedWriter;
import java.io.DataOutputStream;
import java.io.IOException;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLConnection;
import java.net.URLEncoder;

/**
 * SmsClient is the client class for posting newly arrived text messages
 * to the Sms server.
 */
public class SmsClient {
    public SmsClient(String smsServerURL) {
        this.smsServerURL = smsServerURL;
    }

    public boolean postNewSms(SmsMessage newMsg) throws IOException {
        URL url = new URL(smsServerURL);
        HttpURLConnection conn = (HttpURLConnection) url.openConnection();
        conn.setDoOutput(true);
        conn.setRequestMethod("POST");
        OutputStream outputStream = conn.getOutputStream();

        // TODO - Come back here and fix this part after user accounts are properly implemented.
        String paramStr = String.format("userid=%d&origin=%s&msgbody=%s",
                onlyUserID,
                newMsg.getOriginatingAddress(),
                newMsg.getMessageBody());

        byte[] encodedBytes = paramStr.getBytes("UTF-8");
        try (DataOutputStream dos = new DataOutputStream(conn.getOutputStream())) {
            dos.write(encodedBytes);
        }

        int responseCode = conn.getResponseCode();
        if (responseCode == HttpURLConnection.HTTP_OK) {
            Log.d(SmsClient.class.getName(), "Successfully posted new sms");
            return true;
        } else {
            Log.e(SmsClient.class.getName(), "Failed posting new sms " + conn.getResponseMessage());
            return false;
        }
    }

    private String smsServerURL;
    private static final int onlyUserID = 1;
}