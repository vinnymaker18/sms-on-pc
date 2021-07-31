const SERVER_ROOT = "https://minibox.vinayemani.xyz/";

function fetchNewSMS(userID) {
  fetch(SERVER_ROOT + "sms?userid="+userID)
    .then(resp => {
      return resp.json();
    })
    .then(data => {
      let msgIDs = [];
      for (var sms of data) {
        msgIDs[msgIDs.length] = sms.MsgID;
        sendNotification(sms);
      }

      // send a post request marking these smses as read.
      if (msgIDs.length > 0) {
        fetch(SERVER_ROOT + "sms/mark", {
          method: "POST",
          body: JSON.stringify({msgids: msgIDs})
        }).then(resp => {
          return resp.json();
        }).then(_ => {
        }).catch(err => {
          console.log(err);
        });
      }
    })
    .catch(err => {
      console.log(err);
    });
}

function sendNotification(smsMsg) {
  chrome.notifications.create({
    type: "basic",
    iconUrl: "bird.jpeg",
    title: smsMsg.OriginAddress,
    message: smsMsg.MsgBody,
    priority: 2
  });
}

chrome.alarms.create({
  when: 5000,
  periodInMinutes: 0.2
});

chrome.alarms.onAlarm.addListener(() => {
  fetchNewSMS(1);
});
