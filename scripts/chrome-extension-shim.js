/**
 * This is a hacky shim that is injected into the chrome extension. Other
 * extensions aren't yet supported and may never be. I hope everyone is using
 * chrome. This acts as a proxy so that we can take advantage of the fact that
 * the browser extension for chrome is authorized to make calls to 1password
 */

function loadSudolikeabossProxyOnDomReady() {
  var slabServerConnection = new WebSocket('ws://127.0.0.1:16263/browser');
  var onepassServerConnection = new WebSocket('ws://127.0.0.1:6263/4');
  var lastUsedClientId;

  slabServerConnection.onmessage = function(event) {
    var envelope = JSON.parse(event.data);
    lastUsedClientId = envelope.slabClientId;
    onepassServerConnection.send(JSON.stringify(envelope.command));
  };

  slabServerConnection.onerror = function(err) {
    console.error('Some error occured talking to the sudolikeaboss server');
    console.error(err);
  };


  onepassServerConnection.onerror = function(err) {
    console.error('Some error occured talking to the 1password server');
    console.error(err);
  };

  onepassServerConnection.onmessage = function(event) {
    var response = JSON.parse(event.data);
    var envelope = {
      slabClientId: lastUsedClientId,
      response: response
    };
    slabServerConnection.send(JSON.stringify(envelope));
  };
}

document.addEventListener('DOMContentLoaded', loadSudolikeabossProxyOnDomReady, false);
