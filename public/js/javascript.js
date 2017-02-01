var app = angular.module('App', []);

app.controller('AppCtrl', function($scope, $http) {

	$scope.pluginSources=[]
	$scope.componentList=[]
	$scope.functionList=[]
	$scope.state="starting ..."

	$scope.tableResult=[]

	$scope.VisibledTable = false;

	$scope.addDataSource = function() {
		console.log("addDataSource : " + $scope.SelectedDataSource)
		$scope.pluginSources.push($scope.SelectedDataSource)
		$scope.SelectedDataSource = "Choose a plugin" //Reset
		$scope.refreshData();
	}
	$scope.addComponent = function() {
		console.log("addComponent : " + $scope.nameComponent)
		$scope.componentList.push($scope.nameComponent)
		$scope.nameComponent = "" //Reset
		$scope.refreshData();
	}
	$scope.addFunction = function() {
		console.log("addFunction : " + $scope.nameFunction)
		$scope.functionList.push($scope.nameFunction)
		$scope.nameFunction = "" //Reset
		$scope.refreshData();
	}

	$scope.buildQuery = function() {
		return "plugins="+$scope.pluginSources.join()+"&components="+$scope.componentList.join()+"&functions="+$scope.functionList.join()
	}

	$scope.refreshData = function() {
		$scope.state="loading ..."
		$scope.tableResult=[]
		$http.get('../api/list?'+$scope.buildQuery())
			.then(function(response) {
			  $scope.state="parsing results ..."
				console.log("Debug response", response.data)
				for (collector of Object.keys(response.data)) {
					for (soft of response.data[collector]) {
						if (soft.Vulns.length > 0) {
							for (vuln of soft.Vulns) {
								var element = {};
								element.Software = soft.Component.Name;
								element.Version = soft.Component.Version;
								element.Function = soft.Component.Function;
								element.IDVuln = vuln.Value.ID;
								element.Link = vuln.Value.URL;

								$scope.tableResult.push(element)
							}
						}
					}
				}
				$scope.state="ready !"
			});
		/*
		//TODO Find a best way to delete it
		if ($scope.tableResult.length != 0) {
			while($scope.tableResult.length != 0 ) {
				$scope.tableResult.pop()
			}
		}

		for (soft of $scope.Result) {
			if (soft.Vulns.length > 0) {
				for (vuln of soft.Vulns) {
					var element = {};
					element.Software = soft.Component.Name;
					element.Version = soft.Component.Version;
					element.Function = soft.Component.Function;
					element.IDVuln = vuln.Value.ID;
					element.Link = vuln.Value.URL;

					$scope.tableResult.push(element)
				}
			}
		}
		*/
	}

	$scope.loadConfig = function() {
		$http.get('../api/info').
		then(function(response) {
			$scope.info = response.data;
			$scope.info.Functions = []
			for (soft of $scope.info.Components) {
				for (fp of soft.Function.split(",")) {
					if ( $scope.info.Functions.indexOf(fp) == -1 ) {
						$scope.info.Functions.push(fp)
					}
				}
			}
			$scope.state="ready !"
		});
	}
/*
	$scope.loadResults = function() {
		$http.get('../api/list').
		then(function(response) {
			$scope.ListResult = response.data;
			console.log($scope.sources)
			console.log("Results : " + response.data)
		});
	}
*/
	$scope.loadConfig();
});
