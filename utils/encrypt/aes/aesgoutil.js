var crypto;
try {
    crypto = require('crypto');
} catch (err) {
    console.log('crypto support is disabled!');
}

var key = "5a673bd785831e2e5a673bd785831e2e"// crypto.randomBytes(32);
var iv = "5a673bd785831e2e"// crypto.randomBytes(16);

function encrypt(text) {
    // Creating Cipheriv with its parameter
    let cipher =
        crypto.createCipheriv('aes-256-cbc', key, iv);

    // Updating text
    let encrypted = cipher.update(text);

    // Using concatenation
    encrypted = Buffer.concat([encrypted, cipher.final()]);

    // Returning iv and encrypted data
    return {
        iv: iv.toString('hex'),
        key: key.toString('hex'),
        encryptedData: encrypted.toString('hex')
    };
}

console.log(encrypt("hello world"))

// 0c5b82ab4886cbcb6f2dcef07ed276ffff8a87f5525a67bfe49e1eccb5666cb1
// A decrypt function
function decrypt(text, iv) {

    // let iv = Buffer.from(iv, 'hex');
    let encryptedText =
        Buffer.from(text, 'hex');

    // Creating Decipher
    let decipher = crypto.createDecipheriv(
        'aes-256-cbc', Buffer.from(key), iv);

    // Updating encrypted text
    let decrypted = decipher.update(encryptedText);
    decrypted = Buffer.concat([decrypted, decipher.final()]);

    // returns data after decryption
    return decrypted.toString();
}
console.log(decrypt("0c5b82ab4886cbcb6f2dcef07ed276ffff8a87f5525a67bfe49e1eccb5666cb1", iv))