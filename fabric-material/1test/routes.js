
var record = require('./controller.js');

module.exports = function(app){

  app.get('/add_record/:record', function(req, res){
    record.add_record(req, res);
  });
  app.get('/get_all_record', function(req, res){
    record.get_all_record(req, res);
  });
  app.get('/add_policy/:policy', function(req, res){
    record.add_policy(req, res);
  });
  app.get('/flight_details/:flight', function(req, res){
    record.flight_details(req, res);
  });

}
