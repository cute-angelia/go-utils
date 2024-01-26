## readme

js 解密

```js
const {
  logger,
  logerr
} = require('../../libs/logger/logger');
const request = require("../../libs/http/request");
const {
  Secure
} = require("mali-secure");
const crypto = require("crypto");

class RequestService {

  constructor() { }
  static async postRequest(url, token, data) {
    var that = this;
    try {
      const AppId = "daodao";
      const secret = "0db00084fa02e94e2ea5c015d77cdb21";
      const version = "v1.0.1";
      const secure = new Secure(AppId, 1, "", secret, version);
      const headers = {
        'Content-Type': 'application/x-www-form-urlencoded',
        'Authorization': 'Bearer ' + token
      };

      url = url + "?crypto=crypto"
      var response = await request({
        url: secure.getSign(url),
        method: "post",
        data,
        headers
      });
      console.log(response, typeof response.data);
      if (response && response.code == 0) {
        if (response.data == null || typeof response.data === 'object' || Array.isArray(response.data)) {
        } else {
          let key = "adodoaz318x3jvaf" + response.data.substr(0, 16);
          let ciphertext = response.data.substr(16, response.data.length);
          response.data = RequestService.decrypt(ciphertext, key)
          console.log(response); // 'my message'
        }
      }
      return response;
    } catch (error) {
      logger.info('Request failed error=', error);
      logerr.error('Request failed');
    }
  }

  // 解密数据
  static decrypt(cipherText, key) {
    const ALGORITHM = 'aes-256-cbc';
    const BLOCK_SIZE = 16;

    const contents = Buffer.from(cipherText, 'base64');
    const iv = contents.subarray(0, BLOCK_SIZE);
    const textBytes = contents.subarray(BLOCK_SIZE);

    const decipher = crypto.createDecipheriv(ALGORITHM, key, iv);
    let decrypted = decipher.update(textBytes, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    cipherText = JSON.parse(decrypted)
    return cipherText
  }

}
module.exports = RequestService;
```