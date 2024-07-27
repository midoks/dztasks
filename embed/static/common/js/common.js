
function addVersion(){
	console.log("add");


	layer.open({
		type:1,
		title:"插件安装",
		btn: ['确定','关闭'],
		content: "<div class='bt-form pd20 c6'>\
			<div class='version line'></div>\
	    </div>",
	    success:function(){
	    	console.log("ddd");
	    }
	});
}




