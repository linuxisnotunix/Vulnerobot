var app = angular .module('App',[]);
app.controller('AppCtrl',function ($scope){
/*	$scope.tasks = [];
	$scope.query = "";
	$scope.addTask = function(){
		$scope.task.checkbox = false;
		$scope.tasks.push($scope.task);
		$scope.task={};
	}
	$scope.deleteThis = function(t){
		index = $scope.tasks.indexOf(t);
		$scope.tasks.splice(index,1);
	}
	$scope.deleteAllTasksDone= function(){
		for (t of $scope.tasks){
			if (t.checkbox == true) {
				$scope.deleteThis(t);
			}
		}
	}*/
	$scope.query ="";
	$scope.sources = [{"Name":"ANSSI"},{"Name":"NVD"}]
	$scope.SelectedDataSource = $scope.sources[0].Name;

	var ListLogAndFP = [
		{"idUnique":1,"value":"Word","type":"Software"},
		{"idUnique":2,"value":"Excel","type":"Software"},
		{"idUnique":3,"value":"FP1","type":"FP"}
	]

	$scope.ListEnter = ListLogAndFP
	var ANSSIResult=[
		{"Software":"Word","Version":"11.3","Function":"FP1","IDVuln":"CVE343444","NameVuln":"Rights problems","Date":"01/02/16","Link":"link1"},
		{"Software":"Reader","Version":"45A","Function":"FP3","IDVuln":"CVE23444","NameVuln":"Rights problems","Date":"04/02/16","Link":"link2"}
	]
	var NVDResult=[
		{"Software":"databaseReader","Version":"33","Function":"FP2","IDVuln":"CVE365944","NameVuln":"Title 1","Date":"01/03/16","Link":"link3"},
		{"Software":"PDFReader","Version":"46A","Function":"FP4","IDVuln":"CVE23374","NameVuln":"Title 2","Date":"07/09/16","Link":"link4"}
	]
	
	$scope.loadData = function(){
		//$scope.Result='$scope.SelectedDataSource'Result	
		switch($scope.SelectedDataSource){
			case "ANSSI":
				$scope.Result=ANSSIResult
				break
			case "NVD":
				$scope.Result=NVDResult
				break
			default:
				$scope.Result=ANSSIResult

		}
	}

		$scope.loadData();

});