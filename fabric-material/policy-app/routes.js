
var record = require('./controller.js');

module.exports = function(app){

  app.get('/add_record/:record', function(req, res){
    record.add_record(req, res);
  });
  app.get('/get_all_record', function(req, res){
    record.get_all_record(req, res);
  });
  app.get('/select_policy/:policy', function(req, res){
    record.select_policy(req, res);
  });
  app.get('/flight_detail/:flight', function(req, res){
    record.flight_detail(req, res);
  });
  app.get('/roll_claim/', function(req, res){
    record.roll_claim(req, res);
  });
}
