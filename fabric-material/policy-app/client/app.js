
'use strict';

var app = angular.module('application', []);

// Angular Controller
app.controller('appController', function($scope, appFactory){

	$("#success_policy").hide();	
	$("#success_flight").hide();	
	$("#success_roll").hide();
	$("#success_create").hide();
	$("#error_policy").hide();		
	$("#error_flight").hide();		
	$("#error_roll").hide();		

	
	$scope.viewLedger = function(){

		appFactory.viewLedger(function(data){
			var array = [];
			for (var i = 0; i < data.length; i++){
				parseInt(data[i].Key);
				data[i]["Record"].Key = parseInt(data[i].Key);
				array.push(data[i]["Record"]);
			}
			array.sort(function(a, b) {
			    return parseFloat(a.Key) - parseFloat(b.Key);
			});
			$scope.all_record = array;
		});
	}

	$scope.newUser = function(){

		appFactory.newUser($scope.record, function(data){
			$scope.create_record = data;
			$("#success_create").show();
		});
	}

	$scope.selectPolicy = function(){

		appFactory.selectPolicy($scope.policy, function(data){
			$scope.select_policy = data;
			if ($scope.select_policy == "Error: no data found"){
				$("#error_policy").show();
				$("#success_policy").hide();
			} else{
				$("#success_policy").show();
				$("#error_policy").hide();
			}
		});
	}

	$scope.flightDet = function(){

		appFactory.flightDet($scope.flight, function(data){
			$scope.flight_detail = data;
			if ($scope.flight_detail == "Error: no data found"){
				$("#error_flight").show();
				$("#success_flight").hide();
			} else{
				$("#success_flight").show();
				$("#error_flight").hide();
			}
		});
	}

	$scope.rollClaim = function(){

		appFactory.rollClaim(function(data){
			if (data != null){
				$("#success_roll").show();
				$("#error_roll").hide();
			} else{
				$("#success_roll").hide();
				$("#error_roll").show();
			}
		});  
	}


});

// Angular Factory
app.factory('appFactory', function($http){
	
	var factory = {};

    factory.viewLedger = function(callback){

    	$http.get('/get_all_record/').success(function(output){
			callback(output)
		});
	}

	factory.newUser = function(data, callback){


		var record = data.id + "$" + data.userid + "$" + data.name + "$" + data.email + "$" + data.age;

    	$http.get('/add_record/'+record).success(function(output){
			callback(output)
		});
	}

	factory.selectPolicy = function(data, callback){

		var policy = data.id + "$" + data.polid + "$" + data.polname + "$" + data.polval + "$" + data.polvalid;


    	$http.get('/select_policy/'+policy).success(function(output){
			callback(output)
		});
	}

	factory.flightDet = function(data, callback){

		var flight = data.id + "$" + data.flightnum + "$" + data.airline + "$" + data.arrdep + "$" + data.iatacode + "$" + data.time;


    	$http.get('/flight_detail/'+flight).success(function(output){
			callback(output)
		});
	}

    factory.rollClaim = function(callback){

    	$http.get('/roll_claim').success(function(output){
			callback(output)
    	});
	}

	return factory;
});


