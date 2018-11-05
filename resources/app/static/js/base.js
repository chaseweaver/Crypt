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
		astilectron.showOpenDialog({properties: ['openFile', 'singleSelection'],
			title: 'File(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-file', payload: paths[0]});
				document.getElementById('path').innerHTML = paths[0];
		});
	});

	document.getElementById('encrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		key !== "" ? astilectron.sendMessage({name: 'encrypt', payload: key})
			: astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to encrypt.');
			document.getElementById('path').innerHTML = "...";
	});

	document.getElementById('decrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		key !== "" ? astilectron.sendMessage({name: 'decrypt', payload: key})
			: astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to decrypt.');
			document.getElementById('path').innerHTML = "...";
	});
});

function createCopy(cb) {
	cb.checked ? astilectron.sendMessage({name: 'createCopyChecked'})
		: astilectron.sendMessage({name: 'createCopyUnchecked'})
}

function encryptNames(cb) {
	cb.checked ? astilectron.sendMessage({name: 'encryptNamesChecked'})
		: astilectron.sendMessage({name: 'encryptNamesUnchecked'})
}

function keepExtension(cb) {
	cb.checked ? astilectron.sendMessage({name: 'keepExtensionChecked'})
		: astilectron.sendMessage({name: 'keepExtensionUnchecked'})
}

function logOutput(cb) {
	cb.checked ? astilectron.sendMessage({name: 'logOutputChecked'})
		: astilectron.sendMessage({name: 'logOutputUnchecked'})
}