document.addEventListener('astilectron-ready', function() {

	document.getElementById('console-box').innerHTML = 'Crypt is an open-source AES-256 file encryption program.\n\n1. Select File / Directory\n2. Choose options\n3. Enter password\n4. Encrypt / Decrypt\n\ngithub.com/chaseweaver/Crypt\nchase.weaver34@gmail.com\n\n'

	astilectron.onMessage(function(message) {
		document.getElementById('console-box').innerHTML += message.name + '\n\n';
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('toggle-password').addEventListener('click', function() {
		let pwd = document.getElementById('password-box');
    pwd.type === 'password' ? pwd.type = 'text' : pwd.type = 'password';
	});

	document.getElementById('info').addEventListener('click', function() {
		astilectron.showMessageBox({message: 'Crypt is an open-source AES-256 file encryption program written in golang. Find an error? Want to help? Check out the github page <github.com/chaseweaver/Crypt> for more information. Or, how about you email me? <chase.weaver34@gmail.com> Happy hacking.', title: 'Crypt'})
	});

	document.getElementById('file').addEventListener('click', function() {
		astilectron.showOpenDialog({properties: ['openFile', 'multiSelections'],
			title: 'File(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-file', payload: paths});
		});
	});

	document.getElementById('dir').addEventListener('click', function() {
		astilectron.showOpenDialog({properties: ['openDirectory', 'singleSelection'],
			title: 'Dir(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-dir', payload: paths[0]});
		});
	});
	document.getElementById('encrypt').addEventListener('click', function() {
		let key = document.getElementById('password-box').value;
		key !== '' ? astilectron.sendMessage({name: 'encrypt', payload: key})
			: astilectron.showErrorBox('Missing Password!', 'Please enter a password before trying to encrypt.');
	});

	document.getElementById('decrypt').addEventListener('click', function() {
		let key = document.getElementById('password-box').value;
		key !== '' ? astilectron.sendMessage({name: 'decrypt', payload: key})
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