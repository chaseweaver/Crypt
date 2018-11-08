document.addEventListener('astilectron-ready', function() {
	let helpText = 'Crypt is an open-source AES-256 file encryption program.\n\n1. Select File / Directory\n2. Choose options\n3. Enter password\n4. Encrypt / Decrypt\n\ngithub.com/chaseweaver/Crypt\nchase.weaver34@gmail.com\n\n';
	document.getElementById('console-box').innerHTML = helpText;

	let max = 0;
	astilectron.onMessage(function(message) {
		let elem = document.getElementById('bar'); 

		if (message.name === 'count') {
			max = message.payload;
			elem.innerHTML = '0%';
			elem.style.width = 0 + '%';
		} else if (message.name === "progress") {
			let prog = parseInt((message.payload / max) * 100);
			elem.innerHTML = prog + '%'
			elem.style.width = prog + '%'; 

			if (prog >= 100) {
				elem.innerHTML = '100%';
			}
		} else {
			document.getElementById('console-box').innerHTML += message.name + '\n\n';
		}

		if (message.name === 'stop') {
			elem.innerHTML = 'Complete : ' + message.payload;
		}
	});

	document.getElementById('close').addEventListener('click', function() {
		astilectron.sendMessage({name: 'close'});
	});

	document.getElementById('toggle-password').addEventListener('click', function() {
		let pwd = document.getElementById('password-box');
    pwd.type === 'password' ? pwd.type = 'text' : pwd.type = 'password';
	});

	document.getElementById('info').addEventListener('click', function() {
		astilectron.showMessageBox({message: 'Crypt is an open-source AES-256 file encryption program written in golang. Find an error? Want to help? Check out the github page <github.com/chaseweaver/Crypt> for more information. Or, how about you email me? <chase.weaver34@gmail.com> Happy hacking.', title: 'Crypt'});
	});

	document.getElementById('file').addEventListener('click', function() {
		max = 0;
		astilectron.showOpenDialog({properties: ['openFile', 'multiSelections'],
			title: 'File(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-file', payload: paths});
		});
	});

	document.getElementById('dir').addEventListener('click', function() {
		max = 0;
		astilectron.showOpenDialog({properties: ['openDirectory', 'singleSelection'],
			title: 'Dir(s) to Encrypt/Decrypt'}, function(paths) {
				astilectron.sendMessage({name: 'open-dir', payload: paths[0]});
		});
	});

	document.getElementById('log').addEventListener('click', function() {
		let fs = require('fs');
		let data = document.getElementById('console-box').value.replace(helpText, '');
		astilectron.showSaveDialog({title: "Save Log"}, function(fileName) {
			fs.writeFile(fileName, data, (err) => {
        if (err) {
					alert("An error ocurred creating the file "+ err.message)
        }        
        alert("The file has been succesfully saved");
    	});
		})
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
		: astilectron.sendMessage({name: 'createCopyUnchecked'});
}

function encryptNames(cb) {
	cb.checked ? astilectron.sendMessage({name: 'encryptNamesChecked'})
		: astilectron.sendMessage({name: 'encryptNamesUnchecked'});
}

function keepExtension(cb) {
	cb.checked ? astilectron.sendMessage({name: 'keepExtensionChecked'})
		: astilectron.sendMessage({name: 'keepExtensionUnchecked'});
}

function logOutput(cb) {
	cb.checked ? astilectron.sendMessage({name: 'logOutputChecked'})
		: astilectron.sendMessage({name: 'logOutputUnchecked'});
}