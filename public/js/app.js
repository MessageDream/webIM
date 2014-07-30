var lastReceivedStamp = 0;
var isWait = false;
var chatapp=angular.module('chatapp',[]);
chatapp.config (function($interpolateProvider,$httpProvider){
    $interpolateProvider.startSymbol('//');
    $interpolateProvider.endSymbol('//');
    $httpProvider.defaults.transformRequest=function(data){
        var str = [];  
        for(var p in data){  
           str.push(encodeURIComponent(p) + "=" + encodeURIComponent(data[p]));  
        }  
        return str.join("&");    
    };
    $httpProvider.defaults.headers.post = {  
        'Content-Type': 'application/x-www-form-urlencoded'  
    };
});

chatapp.controller('fetch', ['$scope','$interval','$http','$filter', function($scope,$interval,$http,$filter){
    $scope.send=function (){
        var content={
            uname:$scope.uname,
            roomid:$scope.roomid,
            content: $scope.msg
        };
        $http.post('/lp/post',content).success(function(data){
            $scope.msg='';
        });
    };                                                                                                                                                                                  

    $interval(function(){
        if (isWait) return;
        isWait = true;
        $http.get('/lp/fetch',{params:{lastReceived:lastReceivedStamp,roomid:$scope.roomid}}).success(function(data){
            if (data != null) {
                $filter('orderBy')(data,'-Timestamp');
                if(!$scope.data){
                    $scope.data=new Array();
                }
                angular.forEach(data, function(event){
                    switch (event.Type) {
                        case 0: 
                            if (event.User.Name == $scope.uname) {
                                event.Content='You joined the chat room.';
                            } else {
                                 event.Content=event.User.Name +' joined the chat room.';
                            }
                            event.User.Name='';
                            break;
                        case 1: 
                            event.Content=event.User.Name +' left the chat room.';
                            event.User.Name='';
                            break;
                        case 2: 
                            event.User.Name =event.User.Name+':';
                            break;
                    }
                    $scope.data.unshift(event);
                    lastReceivedStamp = event.Timestamp;
                });
            }
            isWait = false;
        }).error(function(err){
            isWait=false;
        });
    },3000);
}]);

chatapp.controller('websocket', ['$scope','$http', function ($scope,$http) {
    $scope.send=function (){
        var content={
            uname:$scope.uname,
            roomid:$scope.roomid,
            content: $scope.msg
        };
        $http.post('/lp/post',content).success(function(data){
            $scope.msg='';
        });
    }; 

    angular.element(document).ready(function(){
        // Create a socket
        var socket = new WebSocket('ws://' + window.location.host + '/ws/join?uname=' + $scope.uname +"&roomid=" + $scope.roomid);
        // Message received on the socket
        socket.onmessage = function (event) {
            var data = angular.fromJson(event.data);
            if (data != null) {
                if(!$scope.data){
                    $scope.data=new Array();
                }
                switch (data.Type) {
                    case 0: 
                        if (data.User.Name == $scope.uname) {
                            data.Content='You joined the chat room.';
                        } else {
                             data.Content=data.User.Name +' joined the chat room.';
                        }
                        data.User.Name='';
                        break;
                    case 1: 
                        data.Content=data.User.Name +' left the chat room.';
                        data.User.Name='';
                        break;
                    case 2: 
                        data.User.Name =data.User.Name+':';
                        break;
                }
                $scope.data.unshift(data);
                $scope.$digest();
            }
        };
    });
}]);
                                                                                                           