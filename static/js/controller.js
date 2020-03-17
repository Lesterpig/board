angular.module('boardApp', [])
.controller('boardCtrl', function($scope, $interval, $http) {

  $scope.categories = {}
  $scope.lastUpdate = ""

  function pool() {
    $http.get('/data').success(function(data) {
      $scope.categories = data.Services
      $scope.lastUpdate = data.LastUpdate
    })
  }

  pool()
  $interval(pool, 30000)

})
