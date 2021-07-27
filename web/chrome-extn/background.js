function fetchNewSMS(userID) {
  // TODO - Implement this.
  fetch("https://minibox.vinayemani.xyz/sms?userid="+userID)
    .then(resp => {
      return resp.json();
    })
    .then(data => {
      for (var sms of data) {
        sendNotification(sms);
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