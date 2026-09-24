#include <stdio.h>
#include <stdlib.h>
typedef struct Node
{
    int data;
    struct Node *next;
} Node;
typedef struct
{
    Node *front;
    Node *rear;
} LinkQueue;
void initQueue(LinkQueue *q)
{
    q->front = NULL;
    q->rear = NULL;
}
int isEmpty(const LinkQueue *q)
{
    return (q->front == NULL);
}
void Push(LinkQueue *q, int data)
{
    Node *newNode = malloc(sizeof(*newNode));
    if (newNode == NULL) {
    fprintf(stderr, "memory allocation failed\n");
    exit(EXIT_FAILURE);
}
    newNode->next = NULL;
    newNode->data = data;
    if (q->rear == NULL)
    {
        q->front = newNode;
        q->rear = newNode;
        return;
    }
    q->rear->next = newNode;
    q->rear = newNode;
}
int Pop(LinkQueue *q, int *out)
{
    if (isEmpty(q))
        return 0;

    Node *tmp = q->front;
    q->front = q->front->next;
    *out = tmp->data;
    free(tmp);
    if (q->front == NULL)
    {
        q->rear = NULL;
    }
    return 1;
}
int main()
{
    LinkQueue q;
    initQueue(&q);
    Push(&q, 10);
    Push(&q, 20);
    int value;
    Pop(&q, &value);
    printf("value:%d\n", value);
    Pop(&q, &value);
    printf("value:%d\n", value);
    Push(&q, 30);
    Pop(&q, &value);
    printf("value:%d\n", value);
    printf("%d\n", isEmpty(&q));
    return 0;
}