document.addEventListener('astilectron-ready', function() {
	astilectron.onMessage(function(message) {
		console.log(message);
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('dev-tools').addEventListener('click', function() {
		astilectron.sendMessage({name: 'dev-tools'});
	});

	document.getElementById('file').addEventListener('click', function() {
		astilectron.showOpenDialog({properties: ['openFile', 'singleSelection'], title: 'File(s) to Encrypt/Decrypt'}, function(paths) {
			astilectron.sendMessage({name: 'open-file', payload: paths[0]});
			document.getElementById('path').innerHTML = paths[0];
		});
	});

	document.getElementById('encrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		if (key !== "") {
			astilectron.sendMessage({name: 'encrypt', payload: key})
		} else {
			astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to encrypt.');
		}
	});

	document.getElementById('decrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		if (key !== "") {
			astilectron.sendMessage({name: 'decrypt', payload: key})
		} else {
			astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to decrypt.');
		}
	});
});

function createCopy(cb) {
	if (cb.checked) {
		astilectron.sendMessage({name: 'createCopyChecked', payload: key})
	} else {
		astilectron.sendMessage({name: 'createCopyUnchecked', payload: key})
	}
}

function encryptNames(cb) {
	if (cb.checked) {
		astilectron.sendMessage({name: 'encryptNamesChecked', payload: key})
	} else {
		astilectron.sendMessage({name: 'encryptNamesUnchecked', payload: key})
	}
}

function keepExtension(cb) {
	if (cb.checked) {
		astilectron.sendMessage({name: 'keepExtensionChecked', payload: key})
	} else {
		astilectron.sendMessage({name: 'keepExtensionUnchecked', payload: key})
	}
}

function logOutput(cb) {
	if (cb.checked) {
		astilectron.sendMessage({name: 'logOutputChecked', payload: key})
	} else {
		astilectron.sendMessage({name: 'logOutputUnchecked', payload: key})
	}
}