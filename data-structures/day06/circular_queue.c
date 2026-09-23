#include<stdio.h>
#define MAX_SIZE 5
typedef struct{
    int data[MAX_SIZE];
    int front;
    int rear;
}SeqQueue;

void initQueue(SeqQueue *q){
    q->front=0;
    q->rear=0;
    return ;
}
int isEmpty(SeqQueue *q){
    return q->front==q->rear;
}
int isFull(SeqQueue *q){
    if((q->rear+1)%MAX_SIZE==q->front)return 1;
    return 0;
}
int enqueue(SeqQueue *q, int value){
    if(isFull(q)) return 0;
    q->data[q->rear]=value;
    q->rear=(q->rear+1)%MAX_SIZE;
    return 1;

}
int dequeue(SeqQueue *q, int *value){
    if(isEmpty(q)) return 0;
    *value=q->data[q->front];
    q->front=(q->front + 1) % MAX_SIZE;
    return 1;
}
int peek(SeqQueue *q, int *value){
    if(isEmpty(q)) return 0;
    *value = q->data[q->front];
    return 1;
}
int main(void){
    SeqQueue q;
    initQueue(&q);
    enqueue(&q,10);
    enqueue(&q,20);
    enqueue(&q,30);
    int value=0;
    dequeue(&q,&value);
    dequeue(&q,&value);
    enqueue(&q,40);
    enqueue(&q,50);
    enqueue(&q,60);
    printf("front=%d\n",q.front);
    printf("rear=%d\n",q.rear);
    printf("full=%d\n",isFull(&q));
    printf("enqueue 70:%d\n",enqueue(&q,70));
    while(dequeue(&q,&value)){
        printf("%d\n",value);
    }

    printf("empty=%d\n",isEmpty(&q));
    

    return 0;
}