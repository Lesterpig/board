angular.module('boardApp', [])
.controller('boardCtrl', function($scope, $interval, $http) {

  $scope.categories = {}

  function pool() {
    $http.get('/data').success(function(data) {
      $scope.categories = data
    })
  }

  pool()
  $interval(pool, 30000)

})
