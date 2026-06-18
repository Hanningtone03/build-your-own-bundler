(function() {
const __modules = {};
const __cache = {};
function __require(id) {
if (__cache[id]) return __cache[id].exports;
const module = { exports: {} };
__cache[id] = module;
__modules[id](module, module.exports, __require);
return module.exports;
}
__modules["math.js"] = function(module, exports, require) {
function add(a, b) {
return a + b;
}
function multiply(a, b) {
return a * b;
}
module.exports.add = add;
module.exports.multiply = multiply;
};
__modules["main.js"] = function(module, exports, require) {
const { add, multiply } = require("math.js");
console.log('Sum:', add(2, 3));
console.log('Product:', multiply(4, 5));
};
__require("main.js");
})();