document.addEventListener('astilectron-ready', function() {

	document.getElementById('console-box').innerHTML = 'Crypt is an open-source AES-256 file encryption program.\n\n1. Select File / Directory\n2. Choose options\n3. Enter password\n4. Encrypt / Decrypt\n\ngithub.com/chaseweaver/Crypt\nchase.weaver34@gmail.com'

	astilectron.onMessage(function(message) {
		console.log(message);
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('console').addEventListener('click', function() {
		astilectron.sendMessage({name: 'console'});
	});

	document.getElementById('info').addEventListener('click', function() {
		astilectron.sendMessage({name: 'info'});
	});

	document.getElementById('file').addEventListener('click', function() {
		astilectron.showOpenDialog({properties: ['openFile', 'singleSelection'],
			title: 'File(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-file', payload: paths[0]});
				document.getElementById('console-box').value += 'Loaded: ' + paths[0] + '\n';
		});
	});

	document.getElementById('encrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		key !== "" ? astilectron.sendMessage({name: 'encrypt', payload: key})
			: astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to encrypt.');
	});

	document.getElementById('decrypt').addEventListener('click', function() {
		let key = document.getElementById('key').value;
		key !== "" ? astilectron.sendMessage({name: 'decrypt', payload: key})
			: astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to decrypt.');
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