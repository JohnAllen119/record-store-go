#include <stdio.h>
#define CAPACITY 5

typedef struct
{
    int data[CAPACITY];
    int front;
    int rear;
} CircularQueue;
void InitQueue(CircularQueue *q)
{
    q->front = 0;
    q->rear = 0;
}
int IsEmpty(CircularQueue *q)
{
    return q->front == q->rear;
}
int IsFull(CircularQueue *q)
{
    return (q->rear + 1) % CAPACITY == q->front;
}
int Push(CircularQueue *q, int value)
{
    if (IsFull(q))
        return 0;
    q->data[q->rear] = value;
    q->rear = (q->rear + 1) % CAPACITY;
    return 1;
}
int Pop(CircularQueue *q, int *out)
{
    if (IsEmpty(q))
        return 0;
    *out = q->data[q->front];
    q->front = (q->front + 1) % CAPACITY;
    return 1;
}
int QueueSize(CircularQueue *q){
    if(q->rear>=q->front){
        return q->rear-q->front;
    }
    return q->rear + CAPACITY - q->front;
}
int main(void)
{
    CircularQueue q;
    InitQueue(&q);
    Push(&q, 10);
    Push(&q, 20);
    Push(&q, 30);
    Push(&q, 40);
    int out = 0;
    Pop(&q, &out);
    Pop(&q, &out);
    Push(&q, 50);
    Push(&q, 60);
    printf("full size: %d\n", QueueSize(&q));
    printf("%d\n",Push(&q, 70));
    Pop(&q, &out);
    printf("%d\n",out);
    Pop(&q, &out);
    printf("%d\n",out);
    Pop(&q, &out);
    printf("%d\n",out);
    Pop(&q, &out);
    printf("%d\n",out);
    printf("%d\n",Pop(&q,&out));
    Push(&q,60);
    printf("%d\n",QueueSize(&q));
    Pop(&q,&out);
    printf("%d\n",QueueSize(&q));
    return 0;
}