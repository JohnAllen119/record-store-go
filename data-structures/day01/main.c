#include<stdio.h>
#include <stdlib.h>
struct Node {
	int data;
	struct Node *next;
};
void swap(int *a,int *b){
    int c;
    c=*a;
    *a=*b;
    *b=c;
    return ;
}
void doubleValues(int *arr, int length){
    for(int i = 0 ; i<length ; i++){
        *(arr+i)*=2;
    }
    return ;
}
int main(){
    /*int x=10,y=20;
    swap(&x,&y);
    printf("x=%d y=%d\n",x,y);
    int values[4]={1,2,3,4};
    doubleValues(values, 4);
    for(int i=0 ; i<4;i++){
        printf("%d ",values[i]);
    }
    printf("\n");*/
    
    /*int *number=malloc(4*sizeof(int));
    if(number==NULL){
        return 1;
    }
    *(number)=10;
    *(number+1)=20;
    *(number+2)=30;
    *(number+3)=40;
    for(int i=0;i<4;i++){
        printf("%d ",number[i]);
    }
    printf("\n");
    free(number);*/
    struct Node *node=malloc(sizeof(struct Node));
    if (node == NULL) {
	return 1;
}
    node->data=10;
    node->next=NULL;
    printf("%d\n",node->data);
    free(node);
    return 0;
}