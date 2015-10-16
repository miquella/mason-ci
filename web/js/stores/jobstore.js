// TODO: make this all es6ish
var EventEmitter = require('events').EventEmitter;

function Store() {
}

Store.prototype = new EventEmitter();

Object.defineProperty(Store.prototype, 'jobs', {
  get: function() {
    var store = this;
    if (this._jobs === undefined) {
      this._jobs = [];
      var xhr = new XMLHttpRequest();
      xhr.open('GET', '/api/jobs/');
      xhr.addEventListener('load', function() {
        if (this.status !== 200) {
          console.error('job index query failed');
          return;
        }

        store._jobs = JSON.parse(this.response);
        store.emit('change');
      });
      xhr.send();
    }

    return this._jobs;
  },
});

Store.prototype.createBuild = function(jobKey) {
  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/api/jobs/'+jobKey+'/builds/');
  xhr.addEventListener('load', function() {
    if (this.status !== 200) {
      console.error('job index query failed');
      return;
    }

    console.log('build started!');
  });
  xhr.send();
};

module.exports = new Store();
